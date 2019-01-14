package local

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/fatih/color"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

var createAction = common.ComponentAction{
	Handler: func(log common.Logger, config common.Configuration, component common.Component) error {
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
	},
}

var removeAction = common.ComponentAction{
	Handler: func(log common.Logger, config common.Configuration, cmp common.Component) error {
		if container, err := getContainer(cmp); err == nil {
			if container.State == "running" {
				if err := stopContainer(cmp); err != nil {
					return err
				}
			}
			fmt.Printf("Removing container '%s' for component '%s'\n", cmp.DockerId, cmp.Name)
			if err := DockerGetClient().ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{}); err != nil {
				return err
			}
			return nil
		}
		return errors.Errorf("Removing container '%s' for component '%s': Not found. Nothing to remove\n", cmp.Name, cmp.DockerId)
	},
}

var startAction = common.ComponentAction{
	Handler: func(log common.Logger, config common.Configuration, cmp common.Component) error {
		if container, err := getContainer(cmp); err == nil {
			fmt.Printf("Starting container '%s' for component '%s'\n", cmp.DockerId, cmp.Name)

			if err := DockerGetClient().ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{}); err != nil {
				return err
			}
			return nil
		}
		return errors.Errorf("Starting container '%s' for component '%s': Not found. Create it first\n", cmp.Name, cmp.DockerId)
	},
}

var stopAction = common.ComponentAction{
	Handler: func(log common.Logger, config common.Configuration, cmp common.Component) error {
		if container, err := getContainer(cmp); err == nil {
			fmt.Printf("Stopping container '%s' for component '%s'\n", cmp.DockerId, cmp.Name)
			if err := DockerGetClient().ContainerStop(context.Background(), container.ID, nil); err != nil {
				return err
			}
			return nil
		}
		return errors.Errorf("Stopping container '%s' for component '%s': Not found found. Nothing to stop\n", cmp.Name, cmp.DockerId)
	},
}

var pullAction = common.ComponentAction{
	Handler: func(log common.Logger, config common.Configuration, cmp common.Component) error {
		fmt.Printf("pulling Image for '%s' (%s) ... ", cmp.Name, cmp.Image)

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
	},
}

var logsAction = common.RawAction{
	Handler: func(log common.Logger, config common.Configuration, args ...string) error {
		follow := false
		if len(args) == 0 {
			return errors.New(fmt.Sprintf("Missing component Name. Available components = %s", common.ComponentNames(config.CurrentProfile().Components)))
		}
		componentId := args[0]
		componentMap := common.ComponentMap(config.CurrentProfile().Components)

		if len(args) > 1 && args[1] == "-f" {
			follow = true
		}

		if component, ok := componentMap[componentId]; ok {
			return dockerPrintLogs(component, follow)
		}
		return errors.New(fmt.Sprintf("Cannot find component '%s'. Available components = %s", componentId, common.ComponentNames(config.CurrentProfile().Components)))
	}}
