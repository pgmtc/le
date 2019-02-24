package builder

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/pgmtc/le/pkg/common"
)

func setUp() (tmpDir string, mockContext common.Context) {
	tmpDir, _ = ioutil.TempDir("", "le-test-mock")
	os.MkdirAll(tmpDir+"/buildtest", os.ModePerm)
	ioutil.WriteFile(tmpDir+"/buildtest/config.yaml", []byte(""+
		"image: test-image\n"+
		"buildroot: "+tmpDir+"/buildtest/\n"+
		"dockerfile: "+tmpDir+"/buildtest/Dockerfile\n"), 0644)
	ioutil.WriteFile(tmpDir+"/buildtest/Dockerfile", []byte("FROM scratch\nADD . ."), 0644)

	os.MkdirAll(tmpDir+"/buildtest-invalid", os.ModePerm)
	ioutil.WriteFile(tmpDir+"/buildtest-invalid/config.yaml", []byte("some:unparsable:rubbish"), 0644)

	config := common.CreateMockConfig([]common.Component{})
	mockContext = common.Context{
		Config: config,
		Log:    &common.StringLogger{},
	}
	return
}

func Test_parseBuildProperties(t *testing.T) {
	tmpDir, _ := setUp()
	buildDir := tmpDir + "/buildtest"
	err, image, buildDir, dockerFile, buildArgs := parseBuildProperties(buildDir)
	expectedImage := "test-image"
	expectedBuildDir := tmpDir + "/buildtest/"
	expectedDockerfile := tmpDir + "/buildtest/Dockerfile"
	expectedBuildArgs := []string{"arg1:value1"}
	if err != nil {
		t.Errorf("Unexpected error returned: %s", err.Error())
	}
	if image != expectedImage {
		t.Errorf("Expected %s, got %s", expectedImage, image)
	}
	if buildDir != expectedBuildDir {
		t.Errorf("Expected %s, got %s", expectedBuildDir, buildDir)
	}
	if dockerFile != expectedDockerfile {
		t.Errorf("Expected %s, got %s", expectedDockerfile, dockerFile)
	}
	if reflect.DeepEqual(expectedBuildArgs, []string{"arg1:value1", "arg2:value2"}) {
		t.Errorf("Expected %s, got %s", expectedBuildArgs, buildArgs)

	}
	// Test error - non existing build dir
	buildDir = tmpDir + "/non-existing"
	err, image, buildDir, dockerFile, buildArgs = parseBuildProperties(buildDir)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Test error - non non-parsable config file
	buildDir = tmpDir + "/buildtest-invalid"
	err, image, buildDir, dockerFile, buildArgs = parseBuildProperties(buildDir)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}

func Test_buildAction(t *testing.T) {
	tmpDir, mockContext := setUp()
	buildContext := tmpDir + "/buildtest"
	spyBuilder := &SpyBuilder{}
	buildAction := getBuildAction(spyBuilder)

	// Try missing specdir parameter
	if err := buildAction.Run(mockContext, "--specdir"); err == nil {
		t.Errorf("Expected error, got nothing")
	}
	// Test non-existing specdir parameter
	if err := buildAction.Run(mockContext, "--nocache"); err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Test success
	expectedImage := "test-image"
	expectedBuildRoot := buildContext + "/"
	expectedDockerFile := buildContext + "/Dockerfile"
	expectedNoCache := true
	expectedCtx := mockContext

	if err := buildAction.Run(mockContext, "--specdir", buildContext, "--nocache"); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if spyBuilder.Spy.Image != expectedImage {
		t.Errorf("Expected %s, got %s", "test-image", spyBuilder.Spy.Image)
	}
	if spyBuilder.Spy.BuildRoot != expectedBuildRoot {
		t.Errorf("Expected %s, got %s", expectedBuildRoot, spyBuilder.Spy.BuildRoot)
	}
	if spyBuilder.Spy.DockerFile != (buildContext + "/Dockerfile") {
		t.Errorf("Expected %s, got %s", expectedDockerFile, spyBuilder.Spy.DockerFile)
	}
	if spyBuilder.Spy.NoCache != expectedNoCache {
		t.Errorf("Expected %t, got %t", expectedNoCache, spyBuilder.Spy.NoCache)
	}
	if spyBuilder.Spy.Ctx != expectedCtx {
		t.Errorf("Expected %v, got %v", expectedCtx, spyBuilder.Spy.Ctx)
	}

	// Test Builder returning error
	spyBuilder.WantErrorMessage = "artificial error"
	err := buildAction.Run(mockContext, "--specdir", buildContext, "--nocache")
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
	if err != nil && err.Error() != spyBuilder.WantErrorMessage {
		t.Errorf("Expected error message %s, got %s", spyBuilder.WantErrorMessage, err.Error())
	}
}
