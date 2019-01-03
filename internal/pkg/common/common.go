package common

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
)

type HandlerArguments struct {
	debug bool
}
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

func ComponentActionHandler(handler func(component Component, handlerArguments HandlerArguments) error, handlerArguments HandlerArguments) func(args []string) error {
	return func(args []string) error {
		if len(args) == 0 {
			return errors.New(fmt.Sprintf("Missing component Name. Available components = %s", ComponentNames()))
		}

		// If all provided, do for all components
		if args[0] == "all" {
			for _, cmp := range GetComponents() {
				err := handler(cmp, handlerArguments)
				if err != nil {
					color.HiBlack(err.Error())
				}
			}
			return nil
		}

		for _, cmpName := range args {
			if component, ok := (ComponentMap())[cmpName]; ok {
				err := handler(component, handlerArguments)
				if err != nil {
					if len(args) > 1 {
						color.HiBlack(err.Error())
					} else {
						return err
					}

				}
			} else {
				if len(args) > 1 { // Single use
					color.HiBlack("Component '%s' has not been found", cmpName)
				} else { // Multiple use
					return errors.New(fmt.Sprintf("Component %s has not been found. Available components = %s", cmpName, ComponentNames()))
				}
			}
		}
		return nil

	}
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
