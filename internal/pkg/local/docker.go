package local

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)


func isRunning(component Component, runningContainers []string) string {
	componentId := component.dockerId
	for _, runningContainer := range runningContainers {
		if (runningContainer == componentId) {
			return "UP"
		}
	}
	return "DOWN"
}

func getContainers() []string {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	containerNames := make([]string, 0)
	for _, container := range containers {
		containerName := container.Names[0][1:len(container.Names[0])]
		containerNames = append(containerNames, containerName)
	}

	return containerNames
}