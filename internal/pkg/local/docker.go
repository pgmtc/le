package local

import (
	"encoding/binary"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
	"os/user"
	"strconv"
)

func dockerGetLogs(component Component, follow bool) error {
	i, err := dockerGetClient().ContainerLogs(context.Background(), component.dockerId, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     follow,
		Tail:       "40",
	})
	if err != nil {
		log.Fatal(err)
	}

	hdr := make([]byte, 8)
	for {
		_, err := i.Read(hdr)
		if err != nil {
			log.Fatal(err)
		}
		var w io.Writer
		switch hdr[0] {
		case 1:
			w = os.Stdout
		default:
			w = os.Stderr
		}
		count := binary.BigEndian.Uint32(hdr[4:])
		dat := make([]byte, count)
		_, err = i.Read(dat)
		fmt.Fprint(w, string(dat))
	}
}

func dockerGetContainers() (map[string]types.Container, error) {
	containers, err := dockerGetClient().ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	containerMap := make(map[string]types.Container)
	for _, container := range containers {
		fmt.Println(container.Names)
		for _, cName := range container.Names {
			containerName := cName[1:len(cName)]
			containerMap[containerName] = container
		}
	}
	return containerMap, nil
}

func startContainer(component Component) error {
	if container, err := getContainer(component); err == nil {
		fmt.Printf("Starting container '%s' for component '%s'\n", component.dockerId, component.name)

		if err := dockerGetClient().ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Starting container '%s' for component '%s': Not found. Create it first\n", component.name, component.dockerId)
}

func stopContainer(component Component) error {
	if container, err := getContainer(component); err == nil {
		fmt.Printf("Stopping container '%s' for component '%s'\n", component.dockerId, component.name)
		if err := dockerGetClient().ContainerStop(context.Background(), container.ID, nil); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Stopping container '%s' for component '%s': Not found found. Nothing to stop\n", component.name, component.dockerId)
}

func removeContainer(component Component) error {
	if container, err := getContainer(component); err == nil {
		if isContainerRunning(getContainerName(container)) {
			if err := stopContainer(component); err != nil {
				return err
			}
		}
		fmt.Printf("Removing container '%s' for component '%s'\n", component.dockerId, component.name)
		if err := dockerGetClient().ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{}); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Removing container '%s' for component '%s': Not found. Nothing to remove\n", component.name, component.dockerId)
}

func createContainer(component Component) error {
	if isContainerRunning(component.dockerId) {
		return errors.New(fmt.Sprintf("Component %s already exist (%s). If you want to recreate, then please stop and remove it first", component.name, component.dockerId))
	}
	fmt.Printf("Creating container '%s' for component '%s': ", component.dockerId, component.name)
	exposePort := strconv.Itoa(component.containerPort)
	mapPort := strconv.Itoa(component.hostPort)
	var exposedPorts nat.PortSet
	var portMap nat.PortMap


	if (component.containerPort > 0 && component.hostPort > 0) {
		exposedPorts = nat.PortSet{nat.Port(exposePort): struct{}{}}
		portMap = map[nat.Port][]nat.PortBinding{nat.Port(exposePort): {{HostIP: "0.0.0.0", HostPort: mapPort}}}
		fmt.Printf(" port %d will be mapped to host port %d : ", component.containerPort, component.hostPort)
	}

	// Mount AWS login credentials
	usr, _ := user.Current()
	dir := usr.HomeDir
	var awsCliMount []mount.Mount
	awsCliMount = append(awsCliMount, mount.Mount{Type: mount.TypeBind,	Source: dir + "/.aws", Target: "/root/.aws"})

	resp, err := dockerGetClient().ContainerCreate(context.Background(), &container.Config {
		Image: component.image,
		Env: component.env,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		PortBindings: portMap,
		Links: component.links,
		Mounts: awsCliMount,
	}, nil, component.dockerId)
	if err != nil {
		return err
	}

	fmt.Println()

	if err := dockerGetClient().ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	if (err != nil) {
		return err
	}

	return nil
}

func getContainer(component Component) (types.Container, error) {
	var nilCont types.Container
	dockerId := component.dockerId
	containerMap, err := dockerGetContainers()
	if (err != nil) {
		return nilCont, err
	}

	if cont, ok := containerMap[dockerId]; ok {
		return cont, nil
	} else {
		return nilCont, errors.New("Can't find the container")
	}
}

func isContainerRunning(containerId string) bool {
	containers, _ := dockerGetContainers()
	if container, ok := containers[containerId]; ok {
		return container.State == "running"
	}
	return false
}

func dockerGetClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}
	return cli
}

func getContainerName(container types.Container) string {
	return container.Names[0][1:len(container.Names[0])]
}

func dockerListImages() error {
	images, err := dockerGetClient().ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return err
	}
	for _, image := range images {
		repoName := image.RepoTags[0]
		fmt.Println(repoName)
	}
	return nil
}

func pullImage(component Component) error{
	fmt.Printf("pulling image for '%s' (%s) ... ", component.name, component.image)
	out, err := dockerGetClient().ImagePull(context.Background(), component.image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	fmt.Printf("done\n")
	return nil
}