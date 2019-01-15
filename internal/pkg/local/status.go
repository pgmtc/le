package local

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"os"
	"strings"
)

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
