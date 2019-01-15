package common

import (
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type fileSystemConfig struct {
	configLocation string
	configFileName string
	currentProfile Profile // profile-xxx.yaml
	config         Config
}

func FileSystemConfig(configLocation string) Configuration {
	return &fileSystemConfig{
		configLocation: configLocation,
		configFileName: "Config.yaml",
	}
}

func (c *fileSystemConfig) LoadProfile(profileName string) (profile Profile, resultErr error) {
	configDir := c.initConfigDir(c.configLocation)
	out := Profile{}

	fileName := path.Join(configDir, "profile-"+profileName+".yaml")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		resultErr = errors.Errorf("profile does not exist, create it first")
		return
	}

	err := YamlUnmarshall(fileName, &out)
	if err != nil {
		resultErr = err
		return
	}

	profile = out
	return
}

func (c *fileSystemConfig) SaveProfile(profileName string, profile Profile) (fileName string, resultErr error) {
	configDir := c.initConfigDir(c.configLocation)
	fileName = path.Join(configDir, "profile-"+profileName+".yaml")
	if err := YamlMarshall(profile, fileName); err != nil {
		resultErr = errors.Errorf("error writing Config file: %s", err.Error())
		return
	}
	return
}

func (c *fileSystemConfig) SaveConfig() (fileName string, resultErr error) {
	//_, err := c.SaveProfile(c.Config.Profile, c.CurrentProfile())
	//if err != nil {
	//	resultErr = errors.Errorf("error saving current profile: %s", err.Error())
	//	return
	//}
	fileName = path.Join(c.initConfigDir(c.configLocation), c.configFileName)
	if err := YamlMarshall(c.config, fileName); err != nil {
		resultErr = errors.Errorf("error writing Config file: %s", err.Error())
		return
	}
	return
}

func (c *fileSystemConfig) LoadConfig() (resultErr error) {
	//fileName := path.Join(c.initConfigDir(), configFileName)
	fileName := path.Join(ParsePath(c.configLocation), c.configFileName)
	if err := YamlUnmarshall(fileName, &c.config); err != nil {
		resultErr = errors.Errorf("error reading Config file %s: %s", fileName, err.Error())
		return
	}
	configProfile, err := c.LoadProfile(c.config.Profile)
	if err != nil {
		resultErr = errors.Errorf("error loading Config's profile: %s", err.Error())
		return
	}
	c.currentProfile = configProfile
	return
}

func (c *fileSystemConfig) GetAvailableProfiles() (profiles []string) {
	configDir := c.initConfigDir(c.configLocation)
	files, _ := filepath.Glob(configDir + "/profile-*.yaml")

	for _, file := range files {
		profileName := strings.TrimPrefix(file, configDir+"/profile-")
		profileName = strings.TrimSuffix(profileName, ".yaml")
		profiles = append(profiles, profileName)
	}
	return
}

func (c *fileSystemConfig) initConfigDir(configLocation string) (configDir string) {
	configDir = ParsePath(configLocation)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		color.Yellow("Config location '%s' does not exist, creating it", configDir)
		if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	return
}

func (c *fileSystemConfig) CurrentProfile() Profile {
	return c.currentProfile
}

func (c *fileSystemConfig) SetProfile(profileName string, profile Profile) {
	c.currentProfile = profile
	c.config.Profile = profileName
}

func (c *fileSystemConfig) Config() Config {
	return c.config
}
