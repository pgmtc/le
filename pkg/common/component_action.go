package common

import (
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

type ComponentAction struct {
	Handler func(ctx Context, cmp Component) error
}

func (a *ComponentAction) Run(ctx Context, args ...string) error {
	if len(args) == 0 {
		return errors.Errorf("Missing component Name. Available components = %s", ComponentNames(ctx.Config.CurrentProfile().Components))
	}

	// If all provided, do for all components
	if args[0] == "all" {
		for _, cmp := range ctx.Config.CurrentProfile().Components {
			err := a.Handler(ctx, cmp)
			if err != nil {
				color.HiBlack(err.Error())
			}
		}
		return nil
	}

	for _, cmpName := range args {
		if component, ok := (ComponentMap(ctx.Config.CurrentProfile().Components))[cmpName]; ok {
			err := a.Handler(ctx, component)
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
				return errors.Errorf("Component %s has not been found. Available components = %s", cmpName, ComponentNames(ctx.Config.CurrentProfile().Components))
			}
		}
	}
	return nil
}
