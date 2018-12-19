package local

import (
	"encoding/binary"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
	"strconv"
)

func dockerGetLogs(containerId string, follow bool) error {
	i, err := dockerGetClient().ContainerLogs(context.Background(), containerId, types.ContainerLogsOptions{
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
		containerName := container.Names[0][1:len(container.Names[0])]
		containerMap[containerName] = container
	}
	return containerMap, nil
}

func dockerStartContainer(container types.Container) error {
	fmt.Println("Starting container ", container.ID[:10], "... ")
	if err := dockerGetClient().ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func dockerStopContainer(container types.Container) error {
	fmt.Println("Stopping container ", container.ID[:10], "... ")
	if err := dockerGetClient().ContainerStop(context.Background(), container.ID, nil); err != nil {
		return err
	}
	return nil
}

func dockerRemoveContainer(container types.Container) error {
	if isContainerRunning(getContainerName(container)) {
		if err := dockerStopContainer(container); err != nil {
			return err
		}
	}
	fmt.Println("Removing container ", container.ID[:10], "... ")
	if err := dockerGetClient().ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{}); err != nil {
		return err
	}
	return nil
}

func createContainer(component Component) error {
	if isContainerRunning(component.dockerId) {
		return errors.New(fmt.Sprintf("Component %s already exist (%s), please stop and remove it first", component.name, component.dockerId))
	}
	exposePort := strconv.Itoa(component.containerPort)
	mapPort := strconv.Itoa(component.hostPort)
	var exposedPorts nat.PortSet
	var portMap nat.PortMap
	if (component.containerPort > 0 && component.hostPort > 0) {
		exposedPorts = nat.PortSet{nat.Port(exposePort): struct{}{}}
		portMap = map[nat.Port][]nat.PortBinding{nat.Port(exposePort): {{HostIP: "0.0.0.0", HostPort: mapPort}}}
		fmt.Printf("Container port %d will be mapped to host port %d\n", component.containerPort, component.hostPort)
	}
	resp, err := dockerGetClient().ContainerCreate(context.Background(), &container.Config {
		Image: component.image,
		Env: component.env,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		PortBindings: portMap,
	}, nil, component.dockerId)
	if err != nil {
		return err
	}
	err = dockerStartContainer(types.Container{ID: resp.ID})
	if (err != nil) {
		return err
	}
	return nil
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