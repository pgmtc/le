package config

import (
	"encoding/json"
	"github.com/fatih/color"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pkg/errors"
)

var createAction = common.RawAction{
	Handler: func(log common.Logger, config common.Configuration, args ...string) error {
		if len(args) < 1 {
			return errors.Errorf("Missing parameters: profileName [sourceProfile], examples:\n" +
				"    orchard config create my-new-profile\n" +
				"    orchard config create my-new-profile some-old-profile")
		}

		profileName := args[0]
		profile := common.DefaultLocalProfile

		if len(args) > 1 {
			copyFromProfile, err := config.LoadProfile(args[1])
			if err != nil {
				return errors.Errorf("Error when loading profile %s: %s", args[1], err.Error())
			}
			profile = copyFromProfile
		}

		fileName, err := config.SaveProfile(profileName, profile)
		if err != nil {
			return errors.Errorf("Error when saving profile: %s", err.Error())
		}

		log.Infof("Successfully saved profile %s to %s", profileName, fileName)
		return nil
	},
}

var initAction = common.RawAction{
	Handler: func(log common.Logger, config common.Configuration, args ...string) error {
		config.SetProfile("default", common.DefaultLocalProfile)

		fileName, err := config.SaveConfig()
		log.Infof("Config written to %s", fileName)

		fileName, err = config.SaveProfile("default", common.DefaultRemoteProfile)
		log.Infof("Profile written to %s", fileName)

		fileName, err = config.SaveProfile("local", common.DefaultLocalProfile)
		log.Infof("Profile written to %s", fileName)

		fileName, err = config.SaveProfile("remote", common.DefaultRemoteProfile)
		log.Infof("Profile written to %s", fileName)

		return err
	},
}

var statusAction = common.RawAction{
	Handler: func(log common.Logger, config common.Configuration, args ...string) error {
		color.HiWhite("Current profile: %s", config.Config().Profile)
		color.HiWhite("Available profiles: %s", config.GetAvailableProfiles())
		if len(args) > 0 && args[0] == "-v" {
			// Verbose output
			s, _ := json.MarshalIndent(config.CurrentProfile().Components, "", "  ")
			color.White("Components: \n%s\n", s)
		} else {
			color.White("Components: (for more verbose output, add '-v' parameter)")
			for i, cmp := range config.CurrentProfile().Components {
				color.White("   %02d | Name: %s, DockerId: %s, Image: %s", i, cmp.Name, cmp.DockerId, cmp.Image)
			}
		}
		return nil
	},
}

var switchAction = common.RawAction{
	Handler: func(log common.Logger, config common.Configuration, args ...string) error {
		if len(args) < 1 {
			return errors.Errorf("Missing parameter: profileName. Example: orchard config switch my-profile")
		}

		profile, err := config.LoadProfile(args[0])
		if err != nil {
			return errors.Errorf("Error when switching profile: %s", err.Error())
		}

		config.SetProfile(args[0], profile)

		configFile, err := config.SaveConfig()
		if err != nil {
			return errors.Errorf("Error when saving config: %s", err.Error())
		}
		log.Infof("Successfully switched profile to %s. Changes written to %s\n", args[0], configFile)
		return nil
	},
}
