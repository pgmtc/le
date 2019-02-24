package common

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var tmpDir string

func setUp(dirSuffix string) fileSystemConfig {
	tmpDir, _ = ioutil.TempDir("", "le-test-Config-mock")
	targetDir := filepath.Join(tmpDir, dirSuffix)
	return fileSystemConfig{
		configLocation: targetDir,
		configFileName: "Config.yaml",
	}

}

func tearDown() {
	os.RemoveAll(tmpDir)
}

func TestFileSystemConfig(t *testing.T) {
	tmpDir, _ = ioutil.TempDir("", "le-test-Config-mock-constructor")
	FileSystemConfig(tmpDir)
}

func TestFileSystemConfig_initConfigDir(t *testing.T) {
	config := setUp(".le-Config")
	defer tearDown()

	cnfDir := config.initConfigDir(config.configLocation)

	if cnfDir != config.configLocation {
		t.Errorf("Expected cnfDir = %s to be the same as tmpDir = %s", cnfDir, config.configLocation)
	}

	if _, err := os.Stat(cnfDir); os.IsNotExist(err) {
		t.Errorf("Expected cnfDir %s to be created, but it does not exist", cnfDir)
	}

	// Test failure
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, nothing returned in recover")
		}
	}()
	_ = config.initConfigDir("K:../../../../../../../../../../../../../../../crap")

}

func TestFileSystemConfig_SaveAndLoadProfile(t *testing.T) {
	config := setUp(".le-Config")
	defer tearDown()

	testProfile := Profile{
		Components: []Component{
			{
				Name: "test-component",
			},
		},
	}

	fileName, err := config.SaveProfile("test", testProfile)
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	expected := config.configLocation + "/profile-test.yaml"
	if fileName != expected {
		t.Errorf("Expected file name to be %s, got %s", expected, fileName)
	}

	// Test save profile - error
	_, err = config.SaveProfile("/../../../../../../../../../../../../../../../../../../../../../../../../../../../../rubbish", testProfile)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Test load profile
	loadedProfile, err := config.LoadProfile("test")
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	if !reflect.DeepEqual(loadedProfile, testProfile) {
		t.Errorf("loadedProfile and testProfile don't match")
	}

	// Test load profile - error
	_, _ = config.SaveProfile("invalid-profile", testProfile)
	invalidProfileFileName := config.configLocation + "/profile-invalid-profile.yaml"
	ioutil.WriteFile(invalidProfileFileName, []byte("\t\tsome-non-yaml-rubbish"), 0644)
	_, err = config.LoadProfile("invalid-profile")
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
	os.Remove(invalidProfileFileName)

	// Test GetProfiles method
	availableProfiles := config.GetAvailableProfiles()
	if len(availableProfiles) != 1 {
		t.Errorf("Expected available profiles length to be 1. Got %d", len(availableProfiles))
	}

	if availableProfiles[0] != "test" {
		t.Errorf("Expected available profile name to be %s, got %s", "test", availableProfiles[0])
	}

	// Test load non-existing profile
	_, err = config.LoadProfile("non-existing-profile")
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}

func Test_fileSystemConfig_SaveConfig(t *testing.T) {
	config := setUp(".le-Config")
	defer tearDown()

	fileName, err := config.SaveConfig()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	expectedFileName := tmpDir + "/.le-Config/Config.yaml"
	if expectedFileName != fileName {
		t.Errorf("Expected file name %s, got %s", expectedFileName, fileName)
	}

	// Test write failure
	config.configFileName = "../../../../../../../../crap"
	fileName, err = config.SaveConfig()
	if err == nil {
		t.Errorf("Expected rror, got nothing")
	}

}

func Test_fileSystemConfig_LoadConfig(t *testing.T) {
	validConfig := setUp(".le-Config")
	defer tearDown()
	validConfig.config.Profile = "default"
	validConfig.SaveConfig()
	validConfig.SaveProfile("default", Profile{})
	err := validConfig.LoadConfig()
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	validConfig.config.Profile = "non-existing"
	validConfig.SaveConfig()
	err = validConfig.LoadConfig()
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Load Config from non-existing location
	tempDir, _ := ioutil.TempDir("", "le-test-Config-mock")
	targetDirTmp := tempDir + "/.le-Config-to-load"
	cnf := fileSystemConfig{
		configLocation: targetDirTmp,
	}

	err = cnf.LoadConfig()
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

}

func Test_fileSystemConfig_CurrentProfile(t *testing.T) {
	testProfile := Profile{
		Components: []Component{
			{Name: "test-component-1"},
			{Name: "test-component-2"},
		},
	}

	config := fileSystemConfig{
		currentProfile: testProfile,
	}

	returned := config.CurrentProfile()
	if !reflect.DeepEqual(returned, testProfile) {
		t.Errorf("Expected returned profile to equal testprofile")
	}

}

func Test_fileSystemConfig_SetProfile(t *testing.T) {
	testProfile := Profile{
		Components: []Component{
			{Name: "test-component-1"},
			{Name: "test-component-2"},
		},
	}
	config := fileSystemConfig{}

	config.SetProfile("test", testProfile)

	if !reflect.DeepEqual(config.currentProfile, testProfile) {
		t.Errorf("Expected returned profile to equal testprofile")
	}
	if config.config.Profile != "test" {
		t.Errorf("Expected profile name in Config to be set to test, got %s", config.config.Profile)
	}
}

func Test_fileSystemConfig_Config(t *testing.T) {
	config := fileSystemConfig{
		config: Config{
			Profile: "test-profile",
		},
	}

	if !reflect.DeepEqual(config.config, config.Config()) {
		t.Errorf("Expected returned Config to equal provided Config")
	}

}
