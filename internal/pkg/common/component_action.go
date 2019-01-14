package common

import (
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

type ComponentAction struct {
	Handler func(log Logger, config Configuration, cmp Component) error
}

func (a *ComponentAction) Run(log Logger, config Configuration, args ...string) error {

	if len(args) == 0 {
		return errors.Errorf("Missing component Name. Available components = %s", ComponentNames(config.CurrentProfile().Components))
	}

	// If all provided, do for all components
	if args[0] == "all" {
		for _, cmp := range config.CurrentProfile().Components {
			err := a.Handler(log, config, cmp)
			if err != nil {
				color.HiBlack(err.Error())
			}
		}
		return nil
	}

	for _, cmpName := range args {
		if component, ok := (ComponentMap(config.CurrentProfile().Components))[cmpName]; ok {
			err := a.Handler(log, config, component)
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
				return errors.Errorf("Component %s has not been found. Available components = %s", cmpName, ComponentNames(config.CurrentProfile().Components))
			}
		}
	}
	return nil
}
