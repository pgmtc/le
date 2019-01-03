package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/pgmtc/orchard-cli/internal/pkg/builder"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pgmtc/orchard-cli/internal/pkg/config"
	"github.com/pgmtc/orchard-cli/internal/pkg/local"
	"github.com/pgmtc/orchard-cli/internal/pkg/source"
	"os"
	"reflect"
)

func init() {
	// Load config
	if err := common.LoadConfig(); err != nil {
		color.HiRed("Error when loading config: %s", err.Error())
		color.HiRed("Try initializing config directory by running '%s config init'", os.Args[0])
	}
}

func main() {
	args := os.Args[1:]

	modules := make(map[string]func(args []string) error)
	modules["local"] = local.Parse
	modules["source"] = source.Parse
	modules["builder"] = builder.Parse
	modules["config"] = config.Parse

	availableModules := reflect.ValueOf(modules).MapKeys()

	if len(args) == 0 {
		color.Blue("Current profile: %s", common.CONFIG.Profile)
		color.Blue("Please provide module")
		color.Blue(fmt.Sprintf(" syntax : %s [module] [action]", os.Args[0]))
		color.Blue(fmt.Sprintf(" example: %s local status", os.Args[0]))
		color.Blue(fmt.Sprintf(" available modules: %s", availableModules))
		os.Exit(1)
	}

	module := args[0]
	handler := modules[module]

	if handler == nil {
		color.Red(fmt.Sprintf(" Module '%s' does not exist. Available modules = %s", module, availableModules))
		os.Exit(1)
	}

	err := handler(args[1:])
	if err != nil {
		color.Red(err.Error())
	} else {
		color.HiGreen("Finished OK (current profile: %s)", common.CONFIG.Profile)
	}
}
