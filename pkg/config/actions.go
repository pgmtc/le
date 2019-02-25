package config

import (
	"encoding/json"
	"github.com/pgmtc/le/pkg/common"
	"github.com/pkg/errors"
)

var createAction = common.RawAction{
	Handler: func(ctx common.Context, args ...string) error {
		config := ctx.Config
		log := ctx.Log
		if len(args) < 1 {
			return errors.Errorf("Missing parameters: profileName [sourceProfile], examples:\n" +
				"    le config create my-new-profile\n" +
				"    le config create my-new-profile some-old-profile")
		}

		profileName := args[0]
		profile := common.DefaultProfile

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
	Handler: func(ctx common.Context, args ...string) error {
		config := ctx.Config
		log := ctx.Log

		config.SetProfile("default", common.DefaultProfile)
		config.SetRepositoryPrefix("https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/")

		fileName, err := config.SaveConfig(false)
		if err != nil {
			return err
		}
		log.Infof("Config written to %s\n", fileName)

		fileName, err = config.SaveProfile("default", common.DefaultProfile)
		log.Infof("Profile written to %s\n", fileName)

		return err
	},
}

var statusAction = common.RawAction{
	Handler: func(ctx common.Context, args ...string) error {
		log := ctx.Log
		config := ctx.Config

		log.Infof("Auth Prefix: %s\n", config.Config().RepositoryPrefix)
		log.Infof("Current profile: %s\n", config.Config().Profile)
		log.Infof("Available profiles: %s\n", config.GetAvailableProfiles())
		if len(args) > 0 && args[0] == "-v" {
			// Verbose output
			s, _ := json.MarshalIndent(config.CurrentProfile().Components, "", "  ")
			log.Infof("Components: \n%s\n", s)
		} else {
			log.Infof("Components: (for more verbose output, add '-v' parameter)")
			for i, cmp := range config.CurrentProfile().Components {
				log.Infof("   %02d | Name: %s, DockerId: %s, Image: %s\n", i, cmp.Name, cmp.DockerId, cmp.Image)
			}
		}
		return nil
	},
}

var switchAction = common.RawAction{
	Handler: func(ctx common.Context, args ...string) error {
		config := ctx.Config
		log := ctx.Log
		if len(args) < 1 {
			return errors.Errorf("Missing parameter: profileName. Example: le config switch my-profile")
		}

		profile, err := config.LoadProfile(args[0])
		if err != nil {
			return errors.Errorf("Error when switching profile: %s", err.Error())
		}

		config.SetProfile(args[0], profile)

		configFile, err := config.SaveConfig(true)
		if err != nil {
			return errors.Errorf("Error when saving config: %s", err.Error())
		}
		log.Infof("Successfully switched profile to %s. Changes written to %s\n", args[0], configFile)
		return nil
	},
}
