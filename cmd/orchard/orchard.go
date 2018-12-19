package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/pgmtc/orchard-cli/internal/pkg/local"
	"github.com/pgmtc/orchard-cli/internal/pkg/source"
	"os"
	"reflect"
)

func main() {
	args := os.Args[1:]

	modules := make(map[string]func(args []string) error)
	modules["local"] = local.Parse
	modules["source"] = source.Parse

	availableModules := reflect.ValueOf(modules).MapKeys()

	if (len(args) == 0) {
		color.Red("Please provide module")
		color.Red(fmt.Sprintf(" %s [module] [action]", os.Args[0]))
		color.Red(fmt.Sprintf(" example: %s local status", os.Args[0]))
		color.Red(fmt.Sprintf(" available modules = %s", availableModules))
		os.Exit(1)
	}

	module := args[0]
	handler := modules[module]

	if handler == nil {
		color.Red(fmt.Sprintf(" Module '%s' does not exist. Available modules = %s", module, availableModules))
		os.Exit(1)
	}

	err := handler(args[1:])
	if (err != nil) {
		color.Red(err.Error())
	} else {
		color.HiGreen("Finished OK")
	}
}
