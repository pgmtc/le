package local

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"os"
	"os/user"
	"strconv"
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

func startContainer(component common.Component, handlerArguments common.HandlerArguments) error {
	if container, err := getContainer(component); err == nil {
		fmt.Printf("Starting container '%s' for component '%s'\n", component.DockerId, component.Name)

		if err := DockerGetClient().ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Starting container '%s' for component '%s': Not found. Create it first\n", component.Name, component.DockerId)
}

func stopContainer(component common.Component, handlerArguments common.HandlerArguments) error {
	if container, err := getContainer(component); err == nil {
		fmt.Printf("Stopping container '%s' for component '%s'\n", component.DockerId, component.Name)
		if err := DockerGetClient().ContainerStop(context.Background(), container.ID, nil); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Stopping container '%s' for component '%s': Not found found. Nothing to stop\n", component.Name, component.DockerId)
}

func removeContainer(component common.Component, handlerArguments common.HandlerArguments) error {
	if container, err := getContainer(component); err == nil {
		if container.State == "running" {
			if err := stopContainer(component, handlerArguments); err != nil {
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

func createContainer(component common.Component, handlerArguments common.HandlerArguments) error {
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
	awsCliMount = append(awsCliMount, mount.Mount{Type: mount.TypeBind, Source: dir + "/.aws", Target: "/root/.aws"})

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

func pullImage(component common.Component, handlerArguments common.HandlerArguments) error {
	fmt.Printf("pulling Image for '%s' (%s) ... ", component.Name, component.Image)
	out, err := DockerGetClient().ImagePull(context.Background(), component.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	fmt.Printf("done\n")
	return nil
}
