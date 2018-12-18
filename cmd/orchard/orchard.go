package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/pgmtc/orchard/internal/pkg/local"
	"github.com/pgmtc/orchard/internal/pkg/source"
	"os"
)

func main() {

	args := os.Args[1:]
	if (len(args) == 0) {
		fmt.Println("Please provide domain, for example:")
		fmt.Printf("    %s local status\n", os.Args[0])
		os.Exit(1)
	}

	domain := args[0]

	actions := make(map[string]func(args []string) error)
	actions["local"] = local.Parse
	actions["source"] = source.Parse

	handler := actions[domain]
	err := handler(args[1:])
	if (err != nil) {
		color.Red(err.Error())
	} else {
		color.HiGreen("Finished OK")
	}
}
