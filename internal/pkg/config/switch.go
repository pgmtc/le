package config

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pkg/errors"
)

type switchAction struct{}

func (switchAction) Run(log common.Logger, args ...string) error {
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
	log.Infof("Successfully switched profile to %s. Changes written to %s", args[0], configFile)
	return nil
}
