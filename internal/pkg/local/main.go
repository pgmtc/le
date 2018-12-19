package local

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"os"
)

func Parse(args []string) error {
	actions := common.MakeActions()
	actions["status"] = status
	actions["stop"] = dockerActionHandler(dockerStopContainer)
	actions["start"] = dockerActionHandler(dockerStartContainer)
	actions["remove"] = dockerActionHandler(dockerRemoveContainer)
	actions["create"] = createContainerHandler
	actions["logs"] = logsHandler(false)
	actions["watch"] = logsHandler(true)
	return common.ParseParams(actions, args)
}

func createContainerHandler(args[] string) error {
	if len(args) == 0 {
		return errors.New(fmt.Sprintf("Missing component name. Available components = %s", componentNames()))
	}
	componentId := args[0]

	if (componentId == "all") {
		for _, cmp := range getComponents() {
			err := createContainer(cmp)
			if (err != nil) {
				return err
			}
		}
		return nil
	}

	componentMap := componentMap()
	if component, ok := componentMap[componentId]; ok {
		return createContainer(component)
	}
	return errors.New(fmt.Sprintf("Cannot find component '%s'. Available components = %s", componentId, componentNames()))
}

func logsHandler(follow bool) func(args[] string) error {
	return func(args[] string) error {
		if len(args) == 0 {
			return errors.New(fmt.Sprintf("Missing component name. Available components = %s", componentNames()))
		}
		componentId := args[0]
		componentMap := componentMap()
		if component, ok := componentMap[componentId]; ok {
			return dockerGetLogs(component.dockerId, follow)
		}
		return errors.New(fmt.Sprintf("Cannot find component '%s'. Available components = %s", componentId, componentNames()))
	}
}

func dockerActionHandler(handler func(container types.Container) error) func(args[] string) error {
	return func (args[] string) error {
		if len(args) == 0 {
			return errors.New(fmt.Sprintf("Missing component name. Available components = %s", componentNames()))
		}
		componentId := args[0]
		if (componentId == "all") {
			containerMap, err := dockerGetContainers()
			if (err != nil) {
				return err
			}

			for _, cmp := range getComponents() {
				var err error
				if container, ok := containerMap[cmp.dockerId]; ok {
					err = handler(container)
					if (err != nil) {
						return err
					}
				} else {
					fmt.Printf("Starting container: Can't start component '%s'. Container '%s'does not exist\n", cmp.name, cmp.dockerId)
				}
			}
			return nil
		}
		componentMap := componentMap()
		if component, ok := componentMap[componentId]; ok {
			dockerId := component.dockerId
			containerMap, err := dockerGetContainers()
			if (err != nil) {
				return err
			}

			if container, ok := containerMap[dockerId]; ok {
				return handler(container)
			} else {
				return errors.New("Can't find the container")
			}
		}
		return errors.New(fmt.Sprintf("Cannot find component '%s'. Available components = %s", componentId, componentNames()))
	}
}

func status(args[] string) error {
	allComponents := getComponents()
	containerMap, err := dockerGetContainers()
	if (err != nil) {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Component", "Docker Container", "Exists", "State", "HTTP"})
	for _, cmp := range allComponents {
		fmt.Printf("\rChecking state of %s ...", cmp.name)
		exists := "NO"
		state := "missing"
		responding := ""
		if container, ok := containerMap[cmp.dockerId]; ok {
			exists = "YES"
			state = container.State
			if state == "running" {
				responding = isResponding(cmp)
			}
		}

		// Some formatting
		switch exists {
			case "YES": exists = color.HiGreenString(exists)
			case "NO": exists = color.RedString(exists)
		}

		switch state {
			case "running": state = color.HiGreenString(state)
			case "exited": state = color.YellowString(state)
			case "missing": state = color.RedString(state)
		}

		switch responding {
			case "200": responding = color.HiGreenString(responding)
			default: responding = color.RedString(responding)
		}


		table.Append([]string{color.YellowString(cmp.name), color.BlueString(cmp.dockerId), color.YellowString(exists), state, responding})
	}

	fmt.Printf("\r")
	table.Render()
	return nil
}
