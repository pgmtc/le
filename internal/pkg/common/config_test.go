package common

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var tmpDir string

func setUp() fileSystemConfig {
	tmpDir, _ = ioutil.TempDir("", "orchard-test-config-mock")
	targetDir := tmpDir + "/.orchard-config"
	return fileSystemConfig{
		configLocation: targetDir,
	}

}

func tearDown() {
	os.RemoveAll(tmpDir)
}

func TestFileSystemConfig(t *testing.T) {
	tmpDir, _ = ioutil.TempDir("", "orchard-test-config-mock-constructor")
	FileSystemConfig(tmpDir)
}

func TestFileSystemConfig_initConfigDir(t *testing.T) {
	config := setUp()
	defer tearDown()

	cnfDir := config.initConfigDir()

	if cnfDir != config.configLocation {
		t.Errorf("Expected cnfDir = %s to be the same as tmpDir = %s", cnfDir, config.configLocation)
	}

	if _, err := os.Stat(cnfDir); os.IsNotExist(err) {
		t.Errorf("Expected cnfDir %s to be created, but it does not exist", cnfDir)
	}
}

func TestFileSystemConfig_SaveAndLoadProfile(t *testing.T) {
	config := setUp()
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

	// Test load profile
	loadedProfile, err := config.LoadProfile("test")
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	if !reflect.DeepEqual(loadedProfile, testProfile) {
		t.Errorf("loadedProfile and testProfile don't match")
	}

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
