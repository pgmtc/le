package local

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"os"
)

func Parse(args []string) error {
	actions := common.MakeActions()
	actions["status"] = status
	actions["stop"] = componentActionHandler(stopContainer)
	actions["start"] = componentActionHandler(startContainer)
	actions["remove"] = componentActionHandler(removeContainer)
	actions["create"] = componentActionHandler(createContainer)
	actions["pull"] = componentActionHandler(pullImage)
	actions["logs"] = logsHandler(dockerPrintLogs, false)
	actions["watch"] = logsHandler(dockerPrintLogs, true)
	return common.ParseParams(actions, args)
}

func logsHandler(handler func(component Component, follow bool) error, follow bool) func(args[] string) error {
	return func(args[] string) error {
		if len(args) == 0 {
			return errors.New(fmt.Sprintf("Missing component name. Available components = %s", componentNames()))
		}
		componentId := args[0]
		componentMap := componentMap()
		if component, ok := componentMap[componentId]; ok {
			return handler(component, follow)
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


		table.Append([]string{color.YellowString(cmp.name), color.HiBlackString(cmp.dockerId), color.YellowString(exists), state, responding})
	}

	fmt.Printf("\r")
	table.Render()
	return nil
}

func componentActionHandler(handler func(component Component) error) func(args []string) error {
	return func(args []string) error {
		if len(args) == 0 {
			return errors.New(fmt.Sprintf("Missing component name. Available components = %s", componentNames()))
		}

		// If all provided, do for all components
		if args[0] == "all" {
			for _, cmp := range getComponents() {
				err := handler(cmp)
				if (err != nil) {
					color.HiBlack(err.Error())
				}
			}
			return nil
		}

		for _, cmpName := range args {
			if component, ok := (componentMap())[cmpName]; ok {
				err := handler(component)
				if (err != nil) {
					if (len(args) > 1) {
						color.HiBlack(err.Error())
					} else {
						return err
					}

				}
			} else {
				if (len(args) > 1) { // Single use
					color.HiBlack("Component '%s' has not been found", cmpName)
				} else { // Multiple use
					return errors.New(fmt.Sprintf("Component %s has not been found. Available components = %s", cmpName, componentNames()))
				}
			}
		}
		return nil

	}
}
