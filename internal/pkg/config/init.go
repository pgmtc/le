package config

import "github.com/pgmtc/orchard-cli/internal/pkg/common"

type initAction struct{}

func (initAction) Run(log common.Logger, config common.Configuration, args ...string) error {
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
}
