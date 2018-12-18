package common

import (
	"errors"
	"fmt"
	"reflect"
)

func MakeActions() map[string]func(args []string) error {
	return make(map[string]func(args []string) error)
}

func ParseParams(actions map[string]func(args []string) error, args []string) error {
	if (len(args) == 0) {
		keys := reflect.ValueOf(actions).MapKeys()
		return errors.New(fmt.Sprintf("Missing action, available actions: %s", keys))
	}
	action := args[0]
	actionArgs := args[1:]

	if (actions[action] != nil) {
		return actions[action](actionArgs)
	} else {
		availableActions := make([]string, 0, len(actions))
		for k := range actions {
			availableActions = append(availableActions, k)
		}
		return errors.New(fmt.Sprintf("Action '%s' does not exist. Available actions = %s", action, availableActions))
	}
	return nil
}
