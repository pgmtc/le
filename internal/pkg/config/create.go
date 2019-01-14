package config

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pkg/errors"
)

type createAction struct{}

func (createAction) Run(log common.Logger, config common.Configuration, args ...string) error {
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
}
