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
		actions[action](actionArgs)
	} else {
		return errors.New(fmt.Sprintf("Action '%s' does not exist", action))
	}
	return nil
}
