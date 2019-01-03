package common

import (
	"io/ioutil"
	"os"
	"os/user"
	"reflect"
	"testing"
)

var (
	dockerTestingHappened bool = false
)

func TestArrContains(t *testing.T) {
	arr := []string{"element1", "element2", "element3"}
	var emptyArr []string

	if !ArrContains(arr, "element1") {
		t.Errorf("Expected true to be returned")
	}
	if ArrContains(arr, "nonExistent") {
		t.Errorf("Expected false to be returned")
	}
	if ArrContains(arr, "") {
		t.Errorf("Expected false to be returned")
	}
	if ArrContains(emptyArr, "element2") {
		t.Errorf("Expected false to be returned")
	}

}

func TestSkipDockerTesting_true(t *testing.T) {
	dockerTestingHappened = false
	origValue := os.Getenv("SKIP_DOCKER_TESTING")
	os.Setenv("SKIP_DOCKER_TESTING", "true")
	defer evalSkipDockerTesting(t, false, origValue)
	SkipDockerTesting(t)
	dockerTestingHappened = true // It should not get here
}

func TestParsePath(t *testing.T) {
	cwd, _ := os.Getwd()
	usr, _ := user.Current()
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test-relative",
			args: args{path: "relative/path.txt"},
			want: cwd + "/relative/path.txt",
		},
		{
			name: "test-absolute",
			args: args{path: "/absolute/relative/path.txt"},
			want: "/absolute/relative/path.txt",
		},
		{
			name: "test-home",
			args: args{path: "~/somedir/path.txt"},
			want: usr.HomeDir + "/somedir/path.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePath(tt.args.path); got != tt.want {
				t.Errorf("ParsePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_marshall(t *testing.T) {
	type TestType struct {
		Name  string
		Value int
	}

	testData := TestType{
		Name:  "Test",
		Value: 10,
	}

	tmpDir, _ := ioutil.TempDir("", "orchard-test-mock")
	defer os.RemoveAll(tmpDir)

	tmpFile := tmpDir + "/marshall-out.yaml"
	junkFile := tmpDir + "/junk.yaml"
	ioutil.WriteFile(junkFile, []byte("junk file content\n"), 0644)

	// Test write success
	err := marshall(testData, tmpFile)
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	// Test write failure (non existing directory)
	err = marshall(testData, "/non-existing"+tmpFile)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

	outData := TestType{}
	err = unmarshall(tmpFile, &outData)
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	if !reflect.DeepEqual(testData, outData) {
		t.Errorf("In and out does not match")
	}

	// Test Junk Content
	err = unmarshall(junkFile, &outData)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Test Non Existing file
	err = unmarshall(junkFile+"-non-existing", &outData)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

}

func TestSkipDockerTesting_false(t *testing.T) {
	dockerTestingHappened = false
	origValue := os.Getenv("SKIP_DOCKER_TESTING")
	os.Unsetenv("SKIP_DOCKER_TESTING")
	defer evalSkipDockerTesting(t, true, origValue)
	SkipDockerTesting(t)
	dockerTestingHappened = true // It should get here
}

func evalSkipDockerTesting(t *testing.T, expectedValue bool, origValue string) {
	if origValue == "" {
		os.Unsetenv("SKIP_DOCKER_TESTING")
	} else {
		os.Setenv("SKIP_DOCKER_TESTING", origValue)
	}

	if expectedValue != dockerTestingHappened {
		t.Errorf("Unexpected dockerTestingHappened value, expected %t got %t", expectedValue, dockerTestingHappened)
	}
}

func Test_initConfigDir(t *testing.T) {
	origConfigLocation := CONFIG_LOCATION
	defer func() {
		CONFIG_LOCATION = origConfigLocation
	}()

	tmpDir, _ := ioutil.TempDir("", "orchard-test-config-mock")
	defer os.RemoveAll(tmpDir)

	targetDir := tmpDir + "/not-yet-existing"
	CONFIG_LOCATION = targetDir

	cnfDir := initConfigDir()

	if cnfDir != targetDir {
		t.Errorf("Expected cnfDir = %s to be the same as tmpDir = %s", cnfDir, targetDir)
	}

	if _, err := os.Stat(cnfDir); os.IsNotExist(err) {
		t.Errorf("Expected cnfDir %s to be created, but it does not exist", cnfDir)
	}
}

func TestSaveAndLoadProfile(t *testing.T) {
	var origConfigLocation, tmpDir string
	var testProfile Profile
	setUp := func() {
		origConfigLocation = CONFIG_LOCATION
		tmpDir, _ = ioutil.TempDir("", "orchard-test-config-mock")
		CONFIG_LOCATION = tmpDir
		testProfile = Profile{
			Components: []Component{
				{
					Name: "test-component",
				},
			},
		}
	}

	rollBack := func() {
		CONFIG_LOCATION = origConfigLocation
		os.RemoveAll(tmpDir)
	}

	setUp()
	defer rollBack()

	fileName, err := SaveProfile("test", testProfile)
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	expected := tmpDir + "/profile-test.yaml"
	if fileName != expected {
		t.Errorf("Expected file name to be %s, got %s", expected, fileName)
	}

	// Test load profile
	loadedProfile, err := LoadProfile("test")
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	if !reflect.DeepEqual(loadedProfile, testProfile) {
		t.Errorf("loadedProfile and testProfile don't match")
	}

	// Test GetProfiles method
	availableProfiles := GetAvailableProfiles()
	if len(availableProfiles) != 1 {
		t.Errorf("Expected available profiles length to be 1. Got %d", len(availableProfiles))
	}

	if availableProfiles[0] != "test" {
		t.Errorf("Expected available profile name to be %s, got %s", "test", availableProfiles[0])
	}

	// Test load non-existing profile
	_, err = LoadProfile("non-existing-profile")
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Test switch profile method
	err = SwitchProfile("test")
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	if !reflect.DeepEqual(CURRENT_PROFILE, testProfile) {
		t.Errorf("Expected CURRENT_PROFILE to equal testProfile")
	}

	// Test switch profile to non-existing
	err = SwitchProfile("non-existing-profile")
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

}

func TestSaveAndLoadConfig(t *testing.T) {
	var origConfigLocation, origProfile, tmpDir string
	setUp := func() {
		origConfigLocation = CONFIG_LOCATION
		origProfile = CONFIG.Profile
		tmpDir, _ = ioutil.TempDir("", "orchard-test-config-mock")
		CONFIG_LOCATION = tmpDir
		CONFIG.Profile = "test-profile"
	}

	rollBack := func() {
		CONFIG_LOCATION = origConfigLocation
		CONFIG.Profile = origProfile
		os.RemoveAll(tmpDir)
	}

	setUp()
	defer rollBack()

	configFile, err := SaveConfig()
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}
	if configFile == "" {
		t.Errorf("Expected config file name not to be empty")
	}

	// Load config
	CONFIG.Profile = "default" // Reset before load
	err = LoadConfig()
	if err != nil {
		t.Errorf("Unexpected error, but got %s", err.Error())
	}
	if CONFIG.Profile != "test-profile" {
		t.Errorf("Expected CONFIG.Profile = %s to equal test-profile", CONFIG.Profile)
	}

	// Test write fail
	CONFIG_FILE = "///"
	configFile, err = SaveConfig()
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
	err = LoadConfig()
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}

func Test_switchCurrentProfile(t *testing.T) {
	newProfile := Profile{
		SourceLocation: "testLocation",
	}
	SwitchCurrentProfile(newProfile)
	if !reflect.DeepEqual(newProfile, GetCurrentProfile()) {
		t.Errorf("Expected newProfile and getCurrentProfile() to match")
	}
}
