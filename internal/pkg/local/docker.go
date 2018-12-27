package local

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"os"
	"os/user"
	"strconv"
)

func dockerPrintLogs(component Component, follow bool) error {
	if container, err := getContainer(component); err == nil {
		options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: follow, Timestamps: false}
		out, err := dockerGetClient().ContainerLogs(context.Background(), container.ID, options)
		if err != nil {
			return err
		}
		io.Copy(os.Stdout, out)
		return nil
	}
	return errors.Errorf("Error when getting container logs for '%s' (%s)\n", component.Name, component.DockerId)
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
		for _, cName := range container.Names {
			containerName := cName[1:len(cName)]
			containerMap[containerName] = container
		}
	}
	return containerMap, nil
}

func startContainer(component Component) error {
	if container, err := getContainer(component); err == nil {
		fmt.Printf("Starting container '%s' for component '%s'\n", component.DockerId, component.Name)

		if err := dockerGetClient().ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Starting container '%s' for component '%s': Not found. Create it first\n", component.Name, component.DockerId)
}

func stopContainer(component Component) error {
	if container, err := getContainer(component); err == nil {
		fmt.Printf("Stopping container '%s' for component '%s'\n", component.DockerId, component.Name)
		if err := dockerGetClient().ContainerStop(context.Background(), container.ID, nil); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Stopping container '%s' for component '%s': Not found found. Nothing to stop\n", component.Name, component.DockerId)
}

func removeContainer(component Component) error {
	if container, err := getContainer(component); err == nil {
		if container.State == "running" {
			if err := stopContainer(component); err != nil {
				return err
			}
		}
		fmt.Printf("Removing container '%s' for component '%s'\n", component.DockerId, component.Name)
		if err := dockerGetClient().ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{}); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Removing container '%s' for component '%s': Not found. Nothing to remove\n", component.Name, component.DockerId)
}

func createContainer(component Component) error {
	if component.Name == "" || component.DockerId == "" || component.Image == "" {
		return errors.New("Missing container Name, DockerId or Image")
	}
	fmt.Println(component.Name)
	fmt.Println(component.DockerId)
	fmt.Println(component.Image)
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
		fmt.Printf(" port %d will be mapped to host port %d : ", component.ContainerPort, component.HostPort)
	}

	// Mount AWS login credentials
	usr, _ := user.Current()
	dir := usr.HomeDir
	var awsCliMount []mount.Mount
	awsCliMount = append(awsCliMount, mount.Mount{Type: mount.TypeBind, Source: dir + "/.aws", Target: "/root/.aws"})

	resp, err := dockerGetClient().ContainerCreate(context.Background(), &container.Config{
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

	if err := dockerGetClient().ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

func getContainer(component Component) (types.Container, error) {
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

func pullImage(component Component) error {
	fmt.Printf("pulling Image for '%s' (%s) ... ", component.Name, component.Image)
	out, err := dockerGetClient().ImagePull(context.Background(), component.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	fmt.Printf("done\n")
	return nil
}
