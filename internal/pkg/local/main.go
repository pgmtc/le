package local

import (
	"fmt"
	"github.com/pgmtc/orchard/internal/pkg/common"
)

func Parse(args []string) error {
	actions := common.MakeActions()
	actions["status"] = status
	return common.ParseParams(actions, args)
}

func status(args[] string) error {
	allComponents := getComponents()
	runningContainers := getContainers()
	fmt.Println(" ----------------------------------------- ")
	fmt.Println("|     Component     |  Docker  |   HTTP   |")
	fmt.Println("|-----------------------------------------|")
	for _, cmp := range allComponents {
		fmt.Printf("|%18.18s |%9.9s |%9.9s |\n", cmp.name, isRunning(cmp, runningContainers), isResponding(cmp))
	}
	fmt.Println(" ----------------------------------------- ")
	return nil
}
