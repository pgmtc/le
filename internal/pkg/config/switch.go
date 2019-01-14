package config

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pkg/errors"
)

type switchAction struct{}

func (switchAction) Run(log common.Logger, config common.Configuration, args ...string) error {
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
}
