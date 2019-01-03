package config

import (
	"encoding/json"
	"github.com/fatih/color"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pkg/errors"
)

func Parse(args []string) error {
	actions := common.MakeActions()
	actions["status"] = status
	actions["init"] = initialize
	actions["create"] = create
	actions["switch"] = switchProfile
	return common.ParseParams(actions, args)
}

func status(args []string) error {
	color.Blue("Current profile: %s", common.CONFIG.Profile)
	color.Blue("Available profiles: %s", common.GetAvailableProfiles())
	if len(args) > 0 && args[0] == "-v" {
		// Verbose output
		s, _ := json.MarshalIndent(common.GetComponents(), "", "  ")
		color.Blue("Components: \n%s\n", s)
	} else {
		color.Blue("Components: (for more verbose output, add '-v' parameter)")
		for i, cmp := range common.GetComponents() {
			color.Blue("   %02d | Name: %s, DockerId: %s, Image: %s", i, cmp.Name, cmp.DockerId, cmp.Image)
		}
	}

	return nil
}

func switchProfile(args []string) error {
	if len(args) < 1 {
		return errors.Errorf("Missing parameter: profileName. Example: orchard config switch my-profile")
	}

	err := common.SwitchProfile(args[0])
	if err != nil {
		return errors.Errorf("Error when switching profile: %s", err.Error())
	}

	configFile, err := common.SaveConfig()
	if err != nil {
		return errors.Errorf("Erorr when saving config: %s", err.Error())
	}
	color.Blue("Successfully switched profile to %s. Changes written to %s", args[0], configFile)
	return nil
}

func initialize(args []string) error {
	common.SwitchCurrentProfile(common.DefaultProfile())
	common.CONFIG.Profile = "default"

	fileName, err := common.SaveConfig()
	color.Blue("Config written to %s", fileName)

	fileName, err = common.SaveProfile("default", common.DefaultProfile())
	color.Blue("Default profile written to %s", fileName)

	return err
}

func create(args []string) error {
	if len(args) < 1 {
		return errors.Errorf("Missing parameters: profileName [sourceProfile], examples:\n" +
			"    orchard config create my-new-profile\n" +
			"    orchard config create my-new-profile some-old-profile")
	}

	profileName := args[0]
	profile := common.DefaultProfile()

	if len(args) > 1 {
		copyFromProfile, err := common.LoadProfile(args[1])
		if err != nil {
			return errors.Errorf("Error when loading profile %s: %s", args[1], err.Error())
		}
		profile = copyFromProfile
	}

	fileName, err := common.SaveProfile(profileName, profile)
	if err != nil {
		return errors.Errorf("Error when saving profile: %s", err.Error())
	}

	color.Blue("Successfully saved profile %s to %s", profileName, fileName)
	return nil
}
