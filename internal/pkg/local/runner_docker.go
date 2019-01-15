package local

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type DockerRunner struct {
}

func (DockerRunner) Status(ctx common.Context, args ...string) error {
	config := ctx.Config
	var verbose bool
	var follow bool
	var followLength int

	if len(args) > 0 && args[0] == "-v" || len(args) > 1 && args[1] == "-v" {
		verbose = true
	}

	// This could be improved - generalized
	if len(args) > 0 && args[0] == "-f" || len(args) > 1 && args[1] == "-f" {
		follow = true
		switch true {
		case len(args) > 1 && args[0] == "-f":
			i, err := strconv.Atoi(args[1])
			if err == nil {
				followLength = i
			}
		case len(args) > 2 && args[1] == "-f":
			i, err := strconv.Atoi(args[2])
			if err == nil {
				followLength = i
			}
		}
		follow = true
	}

	if !follow {
		return printStatus(config.CurrentProfile().Components, verbose, follow, followLength)
	}
	counter := 0
	for {
		printStatus(config.CurrentProfile().Components, verbose, follow, followLength)
		fmt.Println("Orchard local status: ", time.Now().Format("2006-01-02 15:04:05"))
		counter++
		time.Sleep(1 * time.Second)
		if counter == followLength {
			break
		}
	}

	return nil
}

func (DockerRunner) Create(ctx common.Context, cmp common.Component) error {
	log := ctx.Log
	if cmp.Name == "" || cmp.DockerId == "" || cmp.Image == "" {
		return errors.New("Missing container Name, DockerId or Image")
	}

	if _, err := getContainer(cmp); err == nil {
		return errors.Errorf("Component %s already exist (%s). If you want to recreate, then please stop and remove it first\n", cmp.Name, cmp.DockerId)
	}
	log.Infof("Creating container '%s' for component '%s': ", cmp.DockerId, cmp.Name)
	exposePort := strconv.Itoa(cmp.ContainerPort)
	mapPort := strconv.Itoa(cmp.HostPort)
	var exposedPorts nat.PortSet
	var portMap nat.PortMap

	if cmp.ContainerPort > 0 && cmp.HostPort > 0 {
		exposedPorts = nat.PortSet{nat.Port(exposePort): struct{}{}}
		portMap = map[nat.Port][]nat.PortBinding{nat.Port(exposePort): {{HostIP: "0.0.0.0", HostPort: mapPort}}}
		log.Debugf(" port %d will be mapped to host port %d: ", cmp.ContainerPort, cmp.HostPort)
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
		Image:        cmp.Image,
		Env:          cmp.Env,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		PortBindings: portMap,
		Links:        cmp.Links,
		Mounts:       awsCliMount,
	}, nil, cmp.DockerId)
	if err != nil {
		return err
	}

	log.Infof("\n")
	return nil
}

func (DockerRunner) Remove(ctx common.Context, cmp common.Component) error {
	log := ctx.Log
	if container, err := getContainer(cmp); err == nil {
		if container.State == "running" {
			if err := stopContainer(cmp); err != nil {
				return err
			}
		}
		log.Infof("Removing container '%s' for component '%s'\n", cmp.DockerId, cmp.Name)
		if err := DockerGetClient().ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{}); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Removing container '%s' for component '%s': Not found. Nothing to remove\n", cmp.Name, cmp.DockerId)
}

func (DockerRunner) Start(ctx common.Context, cmp common.Component) error {
	log := ctx.Log
	if container, err := getContainer(cmp); err == nil {

		log.Debugf("Starting container '%s' for component '%s'\n", cmp.DockerId, cmp.Name)
		if err := DockerGetClient().ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Starting container '%s' for component '%s': Not found found. Create it first\n", cmp.Name, cmp.DockerId)
}

func (DockerRunner) Stop(ctx common.Context, cmp common.Component) error {
	log := ctx.Log
	if container, err := getContainer(cmp); err == nil {
		log.Debugf("Stopping container '%s' for component '%s'\n", cmp.DockerId, cmp.Name)
		if err := DockerGetClient().ContainerStop(context.Background(), container.ID, nil); err != nil {
			return err
		}
		return nil
	}
	return errors.Errorf("Stopping container '%s' for component '%s': Not found found. Nothing to stop\n", cmp.Name, cmp.DockerId)
}

func (DockerRunner) Pull(ctx common.Context, cmp common.Component) error {
	log := ctx.Log
	log.Infof("pulling Image for '%s' (%s) ... ", cmp.Name, cmp.Image)

	var pullOptions types.ImagePullOptions

	if strings.Contains(cmp.Image, "dkr.ecr.eu-west-1.amazonaws.com") {
		authString, err := getEcrAuth()
		if err != nil {
			return errors.Errorf("problem when obtaining ecr authentication: %s", err.Error())
		}
		pullOptions = types.ImagePullOptions{
			RegistryAuth: authString,
		}
	}

	out, err := DockerGetClient().ImagePull(context.Background(), cmp.Image, pullOptions)
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
		switch true {
		case event.Error != "":
			return errors.Errorf("\nbuild error: %s", event.Error)
		case event.Progress != "" || event.Status != "":
			log.Debugf("\r%s: %s", event.Status, event.Progress)
			if event.ProgressDetail.Current == 0 {
				log.Debugf("\n")
			}
		case strings.TrimSuffix(event.Stream, "\n") != "":
			log.Debugf("%s", event.Stream)
		}

	}

	log.Infof("\n")
	return nil
}

func (DockerRunner) Logs(ctx common.Context, cmp common.Component) error {
	follow := false
	return dockerPrintLogs(cmp, follow)
}
