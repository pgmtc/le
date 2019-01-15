package local

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"os"
	"strconv"
	"strings"
	"time"
)

var statusAction = common.RawAction{
	Handler: func(ctx common.Context, args ...string) error {
		config := ctx.Config
		var verbose bool
		var follow bool
		var followLength int

		if len(args) > 0 && args[0] == "-v" || len(args) > 1 && args[1] == "-v" {
			verbose = true
		}

		// This could be improved - generalized
		if len(args) > 0 && args[0] == "-f" || len(args) > 1 && args[1] == "-f" {
			follow = true
			switch true {
			case len(args) > 1 && args[0] == "-f":
				i, err := strconv.Atoi(args[1])
				if err == nil {
					followLength = i
				}
			case len(args) > 2 && args[1] == "-f":
				i, err := strconv.Atoi(args[2])
				if err == nil {
					followLength = i
				}
			}
			follow = true
		}

		if !follow {
			return printStatus(config.CurrentProfile().Components, verbose, follow, followLength)
		}
		counter := 0
		for {
			printStatus(config.CurrentProfile().Components, verbose, follow, followLength)
			fmt.Println("Orchard local status: ", time.Now().Format("2006-01-02 15:04:05"))
			counter++
			time.Sleep(1 * time.Second)
			if counter == followLength {
				break
			}
		}

		return nil
	},
}

func printStatus(allComponents []common.Component, verbose bool, follow bool, followLength int) error {

	containerMap, err := dockerGetContainers()
	images, err := dockerGetImages()

	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	if verbose {
		table.SetHeader([]string{"Component", "Image (built or pulled)", "Container Exists (created)", "State", "HTTP"})
	} else {
		table.SetHeader([]string{"Component", "Image (built or pulled)", "Container Exists (created)", "HTTP"})
	}

	for _, cmp := range allComponents {
		exists := "NO"
		imageExists := "NO"
		state := "missing"
		responding := ""
		if container, ok := containerMap[cmp.DockerId]; ok {
			exists = "YES"
			state = container.State
			if state == "running" {
				responding, _ = isResponding(cmp)
			}
		}

		if common.ArrContains(images, cmp.Image) {
			imageExists = "YES"
		}

		// Some formatting
		var imageString string = cmp.Image
		if !verbose {
			imageSplit := strings.Split(imageString, "/")
			imageString = imageSplit[len(imageSplit)-1]
		}
		switch imageExists {
		case "YES":
			imageExists = color.HiWhiteString(imageString)
		case "NO":
			imageExists = color.HiBlackString(imageString)
		}

		switch exists {
		case "YES":
			exists = color.HiWhiteString(cmp.DockerId)
		case "NO":
			exists = color.HiBlackString(cmp.DockerId)
		}

		switch state {
		case "running":
			state = color.HiWhiteString(state)
		case "exited":
			state = color.YellowString(state)
		case "missing":
			state = color.HiBlackString(state)
		}

		switch responding {
		case "200":
			responding = color.HiGreenString(responding)
		default:
			responding = color.HiRedString(responding)
		}

		if verbose {
			table.Append([]string{color.YellowString(cmp.Name), imageExists, color.YellowString(exists), state, responding})
		} else {
			table.Append([]string{color.YellowString(cmp.Name), imageExists, color.YellowString(exists), responding})
		}

	}

	if follow {
		print("\033[H\033[2J") // Clear screen
	}
	fmt.Printf("\r")
	table.Render()
	return nil
}
