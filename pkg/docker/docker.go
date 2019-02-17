package docker

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pgmtc/le/pkg/common"
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

func dockerGetImages() (images []string) {
	out, _ := DockerGetClient().ImageList(context.Background(), types.ImageListOptions{})
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

func dockerGetContainers() map[string]types.Container {
	containers, _ := DockerGetClient().ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})

	containerMap := make(map[string]types.Container)
	for _, cont := range containers {
		for _, cName := range cont.Names {
			containerName := cName[1:len(cName)]
			containerMap[containerName] = cont
		}
	}
	return containerMap
}

func startComponent(component common.Component, logger func(format string, a ...interface{})) error {
	if container, err := getContainer(component); err == nil {
		logger("Starting container '%s' for component '%s'\n", component.DockerId, component.Name)

		if err := DockerGetClient().ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Starting container '%s' for component '%s': Not found. Create it first\n", component.Name, component.DockerId)
}

func stopContainer(component common.Component, logger func(format string, a ...interface{})) error {
	if container, err := getContainer(component); err == nil {
		logger("Stopping container '%s' for component '%s'\n", component.DockerId, component.Name)
		if err := DockerGetClient().ContainerStop(context.Background(), container.ID, nil); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Stopping container '%s' for component '%s': Not found found. Nothing to stop\n", component.Name, component.DockerId)
}

func removeComponent(component common.Component, logger func(format string, a ...interface{})) error {
	if container, err := getContainer(component); err == nil {
		if container.State == "running" {
			if err := stopContainer(component, logger); err != nil {
				return err
			}
		}
		logger("Removing container '%s' for component '%s'\n", component.DockerId, component.Name)
		if err := DockerGetClient().ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{}); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Removing container '%s' for component '%s': Not found. Nothing to remove\n", component.Name, component.DockerId)
}

func createContainer(component common.Component, logger func(format string, a ...interface{})) error {
	if component.Name == "" || component.DockerId == "" || component.Image == "" {
		return errors.New("Missing container Name, DockerId or Image")
	}

	if _, err := getContainer(component); err == nil {
		return errors.Errorf("Component %s already exist (%s). If you want to recreate, then please stop and remove it first", component.Name, component.DockerId)
	}
	logger("Creating container '%s' for component '%s': ", component.DockerId, component.Name)
	exposePort := strconv.Itoa(component.ContainerPort)
	mapPort := strconv.Itoa(component.HostPort)
	var exposedPorts nat.PortSet
	var portMap nat.PortMap

	if component.ContainerPort > 0 && component.HostPort > 0 {
		exposedPorts = nat.PortSet{nat.Port(exposePort): struct{}{}}
		portMap = map[nat.Port][]nat.PortBinding{nat.Port(exposePort): {{HostIP: "0.0.0.0", HostPort: mapPort}}}
		logger(" port %d will be mapped to host port %d: ", component.ContainerPort, component.HostPort)
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

	logger("\n")
	return nil
}

func getContainer(component common.Component) (types.Container, error) {
	var nilCont types.Container
	dockerId := component.DockerId
	containerMap := dockerGetContainers()
	if cont, ok := containerMap[dockerId]; ok {
		return cont, nil
	} else {
		return nilCont, errors.New("Can't find the container")
	}
}

func DockerGetClient() *client.Client {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	return cli
}

func removeImage(component common.Component, logger func(format string, a ...interface{})) error {
	logger("removing Image for '%s' (%s) ... \n", component.Name, component.Image)
	name := "docker"
	args := []string{"rmi", component.Image}
	cmd := exec.Command(name, args...)
	if err := cmd.Run(); err != nil {
		return errors.Errorf("Error when removing the image: %s", err.Error())
	}
	return nil
}

func pullImage(component common.Component, logger func(format string, a ...interface{})) error {
	var pullOptions types.ImagePullOptions
	authString, err := getAuthString(component.Repository)
	if err != nil {
		return errors.Errorf("error when obtaining authentication details: %s", err.Error())
	}
	if authString != "" {
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
		}
		switch true {
		case event.Error != "":
			return errors.Errorf("\nbuild error: %s", event.Error)
		case event.Progress != "" || event.Status != "":
			logger(color.MagentaString("\r%s: %s", event.Status, event.Progress))
			if event.ProgressDetail.Current == 0 {
				logger("\n")
			}
		case strings.TrimSuffix(event.Stream, "\n") != "":
			logger(color.MagentaString("%s", event.Stream))
		}

	}

	logger("\n")
	return nil
}

func printStatus(allComponents []common.Component, verbose bool, follow bool, writer io.Writer) error {

	containerMap := dockerGetContainers()
	images := dockerGetImages()
	table := tablewriter.NewWriter(writer)
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
		writer.Write([]byte("\033[H\033[2J")) // Clear screen
	}
	writer.Write([]byte("\r"))
	table.Render()
	return nil
}
