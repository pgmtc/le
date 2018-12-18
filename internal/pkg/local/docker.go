package local

import (
	"encoding/binary"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
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

func dockerGetContainers() []string {
	containers, err := dockerGetClient().ContainerList(context.Background(), types.ContainerListOptions{})
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

func dockerGetClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}
	return cli
}