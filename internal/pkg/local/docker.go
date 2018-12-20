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
		containerName := container.Names[0][1:len(container.Names[0])]
		containerMap[containerName] = container
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
		if container.State == "running" {
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
	if (component.name == "" || component.dockerId == "" || component.image == "") {
		return errors.New("Missing container name, dockerId or image")
	}
	fmt.Println(component.name)
	fmt.Println(component.dockerId)
	fmt.Println(component.image)
	if _, err := getContainer(component); err == nil {
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

func dockerGetClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}
	return cli
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