package docker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

func dockerGetImages() (images []string, returnErr error) {
	out, err := DockerGetClient().ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	for _, img := range out {
		for _, tag := range img.RepoTags {
			images = append(images, tag)
		}
	}
	return
}

func dockerPrintLogs(component common.Component, follow bool) error {
	if container, err := getContainer(component); err == nil {
		options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: follow, Timestamps: false}
		out, err := DockerGetClient().ContainerLogs(context.Background(), container.ID, options)
		if err != nil {
			return err
		}
		io.Copy(os.Stdout, out)
		return nil
	}
	return errors.Errorf("Error when getting container logs for '%s' (%s)\n", component.Name, component.DockerId)
}

func dockerGetContainers() (map[string]types.Container, error) {
	containers, err := DockerGetClient().ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	containerMap := make(map[string]types.Container)
	for _, container := range containers {
		for _, cName := range container.Names {
			containerName := cName[1:len(cName)]
			containerMap[containerName] = container
		}
	}
	return containerMap, nil
}

func startContainer(component common.Component) error {
	if container, err := getContainer(component); err == nil {
		fmt.Printf("Starting container '%s' for component '%s'\n", component.DockerId, component.Name)

		if err := DockerGetClient().ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Starting container '%s' for component '%s': Not found. Create it first\n", component.Name, component.DockerId)
}

func stopContainer(component common.Component) error {
	if container, err := getContainer(component); err == nil {
		fmt.Printf("Stopping container '%s' for component '%s'\n", component.DockerId, component.Name)
		if err := DockerGetClient().ContainerStop(context.Background(), container.ID, nil); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Stopping container '%s' for component '%s': Not found found. Nothing to stop\n", component.Name, component.DockerId)
}

func removeContainer(component common.Component) error {
	if container, err := getContainer(component); err == nil {
		if container.State == "running" {
			if err := stopContainer(component); err != nil {
				return err
			}
		}
		fmt.Printf("Removing container '%s' for component '%s'\n", component.DockerId, component.Name)
		if err := DockerGetClient().ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{}); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Removing container '%s' for component '%s': Not found. Nothing to remove\n", component.Name, component.DockerId)
}

func createContainer(component common.Component) error {
	if component.Name == "" || component.DockerId == "" || component.Image == "" {
		return errors.New("Missing container Name, DockerId or Image")
	}

	if _, err := getContainer(component); err == nil {
		return errors.New(fmt.Sprintf("Component %s already exist (%s). If you want to recreate, then please stop and remove it first", component.Name, component.DockerId))
	}
	fmt.Printf("Creating container '%s' for component '%s': ", component.DockerId, component.Name)
	exposePort := strconv.Itoa(component.ContainerPort)
	mapPort := strconv.Itoa(component.HostPort)
	var exposedPorts nat.PortSet
	var portMap nat.PortMap

	if component.ContainerPort > 0 && component.HostPort > 0 {
		exposedPorts = nat.PortSet{nat.Port(exposePort): struct{}{}}
		portMap = map[nat.Port][]nat.PortBinding{nat.Port(exposePort): {{HostIP: "0.0.0.0", HostPort: mapPort}}}
		fmt.Printf(" port %d will be mapped to host port %d: ", component.ContainerPort, component.HostPort)
	}

	// Mount AWS login credentials
	usr, _ := user.Current()
	dir := usr.HomeDir
	var awsCliMount []mount.Mount
	awsCliPath := filepath.Join(dir, ".aws")
	if _, err := os.Stat(awsCliPath); !os.IsNotExist(err) {
		awsCliMount = append(awsCliMount, mount.Mount{Type: mount.TypeBind, Source: dir + "/.aws", Target: "/root/.aws"})
	}

	_, err := DockerGetClient().ContainerCreate(context.Background(), &container.Config{
		Image:        component.Image,
		Env:          component.Env,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		PortBindings: portMap,
		Links:        component.Links,
		Mounts:       awsCliMount,
	}, nil, component.DockerId)
	if err != nil {
		return err
	}

	fmt.Println()
	return nil
}

func getContainer(component common.Component) (types.Container, error) {
	var nilCont types.Container
	dockerId := component.DockerId
	containerMap, err := dockerGetContainers()
	if err != nil {
		return nilCont, err
	}

	if cont, ok := containerMap[dockerId]; ok {
		return cont, nil
	} else {
		return nilCont, errors.New("Can't find the container")
	}
}

func DockerGetClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}
	return cli
}

func removeImage(component common.Component) error {
	fmt.Printf("removing Image for '%s' (%s) ... \n", component.Name, component.Image)
	name := "docker"
	args := []string{"rmi", component.Image}
	cmd := exec.Command(name, args...)
	if err := cmd.Run(); err != nil {
		return errors.Errorf("Error when removing the image: %s", err.Error())
	}
	return nil
}

func pullImage(component common.Component) error {
	var pullOptions types.ImagePullOptions

	if strings.Contains(component.Image, "dkr.ecr.eu-west-1.amazonaws.com") {
		authString, err := getEcrAuth()
		if err != nil {
			return errors.Errorf("problem when obtaining ecr authentication: %s", err.Error())
		}
		pullOptions = types.ImagePullOptions{
			RegistryAuth: authString,
		}
	}

	out, err := DockerGetClient().ImagePull(context.Background(), component.Image, pullOptions)
	if err != nil {
		return err
	}
	defer out.Close()

	d := json.NewDecoder(out)

	type Event struct {
		Stream         string `json:"stream"`
		Status         string `json:"status"`
		Error          string `json:"error"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"progressDetail"`
	}

	var event *Event
	for {
		if err := d.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		//fmt.Printf("%+v\n", event)
		switch true {
		case event.Error != "":
			return errors.Errorf("\nbuild error: %s", event.Error)
		case event.Progress != "" || event.Status != "":
			fmt.Printf(color.MagentaString("\r%s: %s", event.Status, event.Progress))
			if event.ProgressDetail.Current == 0 {
				fmt.Println()
			}
		case strings.TrimSuffix(event.Stream, "\n") != "":
			fmt.Printf(color.MagentaString("%s", event.Stream))
		}

	}

	//io.Copy(os.Stdout, out)
	//fmt.Printf("done\n")
	fmt.Printf("\n")
	return nil
}

func getEcrAuth() (authString string, resultError error) {
	name := "aws"
	args := []string{"ecr", "get-login", "--no-include-email", "--region", "eu-west-1"}
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		resultError = errors.Errorf("Error when pulling the image: %s", err.Error())
		return
	}

	authString, resultError = parseAwsLogin(string(out))
	return
}

func parseAwsLogin(loginOutput string) (authString string, resultError error) {
	split := strings.Split(loginOutput, " ")
	if len(split) != 7 {
		resultError = errors.Errorf("Unexpected number of items in aws docker login command, got %d, expected %d", len(split), 7)
		return
	}

	authConfig := types.AuthConfig{
		Username:      split[3],
		Password:      split[5],
		ServerAddress: split[6],
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		resultError = errors.Errorf("Error when encoding ECR auth: %s", err.Error())
	}

	authString = base64.URLEncoding.EncodeToString(encodedJSON)
	return
}

func printStatus(allComponents []common.Component, verbose bool, follow bool, followLength int) error {

	containerMap, err := dockerGetContainers()
	images, err := dockerGetImages()

	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Component", "Image (built or pulled)", "Container Exists (created)", "State", "HTTP"})

	for _, cmp := range allComponents {
		exists := "NO"
		imageExists := "NO"
		state := "missing"
		responding := ""
		if container, ok := containerMap[cmp.DockerId]; ok {
			exists = "YES"
			state = container.State
			if state == "running" {
				responding, _ = isResponding(cmp)
			}
		}

		if common.ArrContains(images, cmp.Image) {
			imageExists = "YES"
		}

		// Some formatting
		var imageString = cmp.Image
		if !verbose {
			imageSplit := strings.Split(imageString, "/")
			imageString = imageSplit[len(imageSplit)-1]
		}
		switch imageExists {
		case "YES":
			imageExists = color.HiWhiteString(imageString)
		case "NO":
			imageExists = color.MagentaString(imageString)
		}

		switch exists {
		case "YES":
			exists = color.HiWhiteString(cmp.DockerId)
		case "NO":
			exists = color.MagentaString(cmp.DockerId)
		}

		switch state {
		case "running":
			state = color.HiWhiteString(state)
		case "exited":
			state = color.WhiteString(state)
		case "missing":
			state = color.MagentaString(state)
		}

		switch responding {
		case "200":
			responding = color.HiGreenString(responding)
		default:
			responding = color.MagentaString(responding)
		}

		table.Append([]string{color.HiWhiteString(cmp.Name), imageExists, color.HiWhiteString(exists), state, responding})

	}

	if follow {
		print("\033[H\033[2J") // Clear screen
	}
	fmt.Printf("\r")
	table.Render()
	return nil
}
