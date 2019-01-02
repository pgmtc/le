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
	actions["stop"] = common.ComponentActionHandler(stopContainer, common.HandlerArguments{})
	actions["start"] = common.ComponentActionHandler(startContainer, common.HandlerArguments{})
	actions["remove"] = common.ComponentActionHandler(removeContainer, common.HandlerArguments{})
	actions["create"] = common.ComponentActionHandler(createContainer, common.HandlerArguments{})
	actions["pull"] = common.ComponentActionHandler(pullImage, common.HandlerArguments{})
	actions["logs"] = logsHandler(dockerPrintLogs, false)
	actions["watch"] = logsHandler(dockerPrintLogs, true)
	return common.ParseParams(actions, args)
}

func logsHandler(handler func(component common.Component, follow bool) error, follow bool) func(args []string) error {
	return func(args []string) error {
		if len(args) == 0 {
			return errors.New(fmt.Sprintf("Missing component Name. Available components = %s", common.ComponentNames()))
		}
		componentId := args[0]
		componentMap := common.ComponentMap()
		if component, ok := componentMap[componentId]; ok {
			return handler(component, follow)
		}
		return errors.New(fmt.Sprintf("Cannot find component '%s'. Available components = %s", componentId, common.ComponentNames()))
	}
}

func status(args []string) error {
	allComponents := common.GetComponents()
	containerMap, err := dockerGetContainers()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Component", "Docker Container", "Exists", "State", "HTTP"})
	for _, cmp := range allComponents {
		fmt.Printf("\rChecking state of %s ...", cmp.Name)
		exists := "NO"
		state := "missing"
		responding := ""
		if container, ok := containerMap[cmp.DockerId]; ok {
			exists = "YES"
			state = container.State
			if state == "running" {
				responding = isResponding(cmp)
			}
		}

		// Some formatting
		switch exists {
		case "YES":
			exists = color.HiGreenString(exists)
		case "NO":
			exists = color.RedString(exists)
		}

		switch state {
		case "running":
			state = color.HiGreenString(state)
		case "exited":
			state = color.YellowString(state)
		case "missing":
			state = color.RedString(state)
		}

		switch responding {
		case "200":
			responding = color.HiGreenString(responding)
		default:
			responding = color.RedString(responding)
		}

		table.Append([]string{color.YellowString(cmp.Name), color.HiBlackString(cmp.DockerId), color.YellowString(exists), state, responding})
	}

	fmt.Printf("\r")
	table.Render()
	return nil
}
