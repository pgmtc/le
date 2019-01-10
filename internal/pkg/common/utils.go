package common

import (
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func ArrContains(arr []string, value string) bool {
	for _, element := range arr {
		if element == value {
			return true
		}
	}
	return false
}

func SkipDockerTesting(t *testing.T) {
	if os.Getenv("SKIP_DOCKER_TESTING") != "" {
		t.Skip("Skipping docker testing")
	}
}

/* Method replaces relative path with absolute and replace ~ with user's home dir */
func ParsePath(path string) (result string) {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	result = strings.Replace(path, "~", usr.HomeDir, 1)

	if !strings.HasPrefix(result, "/") {
		currentDir, _ := os.Getwd()
		result = currentDir + "/" + result
	}
	return
}

func initConfigDir() (configDir string) {
	configDir = ParsePath(CONFIG_LOCATION)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		color.Yellow("Config location '%s' does not exist, creating it", configDir)
		if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	return
}

func SaveProfile(profileName string, profile Profile) (fileName string, resultErr error) {
	configDir := initConfigDir()
	fileName = path.Join(configDir, "profile-"+profileName+".yaml")
	if err := marshall(profile, fileName); err != nil {
		resultErr = errors.Errorf("error writing config file: %s", err.Error())
		return
	}
	return
}

func LoadProfile(profileName string) (profile Profile, resultErr error) {
	configDir := initConfigDir()
	out := Profile{}

	fileName := path.Join(configDir, "profile-"+profileName+".yaml")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		resultErr = errors.Errorf("profile does not exist, create it first")
		return
	}

	err := unmarshall(fileName, &out)
	if err != nil {
		resultErr = err
		return
	}

	profile = out
	return
}

func SwitchCurrentProfile(profile Profile) {
	CURRENT_PROFILE = profile
}

func GetCurrentProfile() Profile {
	return CURRENT_PROFILE
}

func SaveConfig() (fileName string, resultErr error) {
	// Save profile
	_, err := SaveProfile(CONFIG.Profile, CURRENT_PROFILE)
	if err != nil {
		resultErr = errors.Errorf("error saving current profile: %s", err.Error())
		return
	}

	fileName = path.Join(initConfigDir(), CONFIG_FILE)
	if err := marshall(CONFIG, fileName); err != nil {
		resultErr = errors.Errorf("error writing config file: %s", err.Error())
		return
	}
	return
}

func LoadConfig() (resultErr error) {
	fileName := path.Join(initConfigDir(), CONFIG_FILE)
	if err := unmarshall(fileName, &CONFIG); err != nil {
		resultErr = errors.Errorf("error writing config file: %s", err.Error())
		return
	}
	configProfile, err := LoadProfile(CONFIG.Profile)
	if err != nil {
		resultErr = errors.Errorf("error loading config's profile: %s", err.Error())
		return
	}
	SwitchCurrentProfile(configProfile)
	return
}

func marshall(data interface{}, fileName string) (resultErr error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		resultErr = errors.Errorf("error when marshalling config: %s", err.Error())
		return
	}

	if err := ioutil.WriteFile(fileName, bytes, 0644); err != nil {
		resultErr = errors.Errorf("error writing file: %s", err.Error())
		return
	}

	return
}

func SwitchProfile(profileName string) (resultErr error) {
	profile, err := LoadProfile(profileName)
	if err != nil {
		resultErr = errors.Errorf("error when switching profile: %s", err.Error())
		return
	}
	CONFIG.Profile = profileName
	CURRENT_PROFILE = profile

	return
}

func GetAvailableProfiles() (profiles []string) {
	configDir := initConfigDir()
	files, returnError := filepath.Glob(configDir + "/profile-*.yaml")
	if returnError != nil {
		return
	}

	for _, file := range files {
		profileName := strings.TrimPrefix(file, configDir+"/profile-")
		profileName = strings.TrimSuffix(profileName, ".yaml")
		profiles = append(profiles, profileName)
	}
	return
}

func unmarshall(fileName string, out interface{}) (resultErr error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		resultErr = errors.Errorf("error when opening file %s: %s", fileName, err.Error())
		return
	}

	if err := yaml.Unmarshal(bytes, out); err != nil {
		resultErr = errors.Errorf("error when unmarshalling file %s: %s", fileName, err.Error())
		return
	}
	return
}
