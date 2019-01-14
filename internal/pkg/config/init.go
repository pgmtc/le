package config

import "github.com/pgmtc/orchard-cli/internal/pkg/common"

type initAction struct{}

func (initAction) Run(log common.Logger, args ...string) error {
	common.SwitchCurrentProfile(common.DefaultLocalProfile())
	common.CONFIG.Profile = "default"

	fileName, err := common.SaveConfig()
	log.Infof("Config written to %s", fileName)

	fileName, err = common.SaveProfile("default", common.DefaultRemoteProfile())
	log.Infof("Profile written to %s", fileName)

	fileName, err = common.SaveProfile("local", common.DefaultLocalProfile())
	log.Infof("Profile written to %s", fileName)

	fileName, err = common.SaveProfile("remote", common.DefaultRemoteProfile())
	log.Infof("Profile written to %s", fileName)

	return err
}
