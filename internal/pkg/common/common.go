package common

import (
	"errors"
	"fmt"
)

func MakeActions() map[string]func(args []string) error {
	return make(map[string]func(args []string) error)
}

func ParseParams(actions map[string]func(args []string) error, args []string) error {
	if len(args) == 0 {
		return errors.New(fmt.Sprintf("Missing action, available actions = %s", getActionNames(actions)))
	}
	if actions["help"] == nil { // Add generic help handler
		actions["help"] = mkHelpHandler(actions)
	}
	action := args[0]
	actionArgs := args[1:]
	if actions[action] != nil {
		return actions[action](actionArgs)
	}
	return errors.New(fmt.Sprintf("Action '%s' does not exist. Available actions = %s", action, getActionNames(actions)))
}

func getActionNames(actions map[string]func(args []string) error) []string {
	availableActions := make([]string, 0, len(actions))
	for k := range actions {
		availableActions = append(availableActions, k)
	}
	return availableActions
}

func mkHelpHandler(actions map[string]func(args []string) error) func(args []string) error {
	return func(args []string) error {
		fmt.Printf("Available actions = %s\n", getActionNames(actions))
		return nil
	}
}
