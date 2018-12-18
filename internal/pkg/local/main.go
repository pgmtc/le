package local

import (
	"errors"
	"fmt"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

func Parse(args []string) error {
	actions := common.MakeActions()
	actions["status"] = status
	actions["logs"] = logsHandler(false)
	actions["watch"] = logsHandler(true)
	return common.ParseParams(actions, args)
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

func status(args[] string) error {
	allComponents := getComponents()
	runningContainers := dockerGetContainers()
	fmt.Println(" ----------------------------------------- ")
	fmt.Println("|     Component     |  Docker  |   HTTP   |")
	fmt.Println("|-----------------------------------------|")
	for _, cmp := range allComponents {
		fmt.Printf("|%18.18s |%9.9s |%9.9s |\n", cmp.name, isRunning(cmp, runningContainers), isResponding(cmp))
	}
	fmt.Println(" ----------------------------------------- ")

	return nil
}
