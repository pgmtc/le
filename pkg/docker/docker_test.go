package docker

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/pgmtc/le/pkg/common"
)

func setUpForRun() (logger common.Logger) {
	logger = &common.StringLogger{}
	return
}

func setUpForBuild() (tmpDir string, mockContext common.Context) {
	tmpDir, _ = ioutil.TempDir("", "le-test-mock")
	os.MkdirAll(tmpDir+"/src", os.ModePerm)
	os.MkdirAll(tmpDir+"/dest", os.ModePerm)

	os.MkdirAll(tmpDir+"/src/subdir", os.ModePerm)
	os.MkdirAll(tmpDir+"/src/.hiddendir", os.ModePerm)

	fileContent := []byte("testing file contents\n")
	ioutil.WriteFile(tmpDir+"/src/"+"file1.txt", fileContent, 0644)
	ioutil.WriteFile(tmpDir+"/src/subdir/"+"file2.txt", fileContent, 0644)
	ioutil.WriteFile(tmpDir+"/src/.hiddendir/"+"file3.txt", fileContent, 0644)
	ioutil.WriteFile(tmpDir+"/src/Dockerfile", []byte("FROM scratch"), 0644)

	os.MkdirAll(tmpDir+"/buildtest", os.ModePerm)
	ioutil.WriteFile(tmpDir+"/buildtest/Dockerfile", []byte("FROM scratch\nADD . ."), 0644)

	config := common.CreateMockConfig([]common.Component{})
	mockContext = common.Context{
		Config: config,
		Log:    &common.StringLogger{},
	}
	return
}

func TestMissingParameters(t *testing.T) {
	logger := setUpForRun()
	cmp := common.Component{
		Name:     "test",
		DockerId: "test-container",
	}
	err := createContainer(cmp, logger.Infof)
	if err == nil {
		t.Errorf("Expected to fail due to mandatory missing")
	}
}

func Test_pullImage(t *testing.T) {
	logger := setUpForRun()
	if os.Getenv("NO_NETWORK") == "true" {
		t.Skipf("NO_NETWORK set to true, skipping")
	}
	common.SkipDockerTesting(t)
	type args struct {
		component common.Component
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "testPublicImage",
			args: args{
				component: common.Component{
					Name:  "test-container",
					Image: "docker.io/library/nginx:stable-alpine",
				},
			},
			wantErr: false,
		},
		{
			name: "testECRWithLogin",
			args: args{
				component: common.Component{
					Name:       "local-db",
					Image:      "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-valuation-client-ui:latest",
					Repository: "ecr:eu-west-1",
				},
			},
			wantErr: !(os.Getenv("SKIP_AWS_TESTING") == ""),
		},
		{
			name: "testNonExistingRepository",
			args: args{
				component: common.Component{
					Name:       "test-container",
					Image:      "non-existing-image",
					Repository: "non-existing",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			removeImage(tt.args.component, logger.Infof) // Ignore errors, image may not exist

			err := pullImage(tt.args.component, logger.Infof)
			if (err != nil) != tt.wantErr {
				t.Errorf("pullImage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				// Check that the image exists
				images := getImages()
				imageClensed := strings.Replace(tt.args.component.Image, "docker.io/library/", "", 1) // This is a workaround for public docker images from docker.io
				if !common.ArrContains(images, imageClensed) {
					t.Logf("%s", images)
					t.Logf("%s", tt.args.component.Image)
					t.Logf("%s", imageClensed)
					t.Errorf("Pulled image '%s' seems not to exist", tt.args.component.Image)
				}

				// Try to remove image
				err = removeImage(tt.args.component, logger.Infof)
				if err != nil {
					t.Errorf("Unexpected error when removing image: %s", err.Error())
				}
				images = getImages()
				if common.ArrContains(images, tt.args.component.Image) {
					t.Errorf("Pulled image '%s' still exist, should have been removed", tt.args.component.Image)
				}
			}

		})
	}
}

func TestComplex(t *testing.T) {
	logger := setUpForRun()
	if os.Getenv("NO_NETWORK") == "true" {
		t.Skipf("NO_NETWORK set to true, skipping")
	}
	common.SkipDockerTesting(t)
	var err error

	cmp1 := common.Component{
		Name:     "linkedContainer",
		DockerId: "linkedContainer",
		Image:    "docker.io/library/nginx:stable-alpine",
	}

	err = pullImage(cmp1, logger.Infof)
	if err != nil {
		t.Errorf("Error, expected Image to be pulled, got %s", err.Error())
	}

	err = createContainer(cmp1, logger.Infof)
	defer removeComponent(cmp1, logger.Infof)
	if err != nil {
		t.Errorf("Error, expected container to be created, got %s", err.Error())
	}
	err = startComponent(cmp1, logger.Infof)
	if err != nil {
		t.Errorf("Error, expected container to be created, got %s", err.Error())
	}

	cmp := common.Component{
		Name:          "test",
		DockerId:      "testContainer",
		Image:         "nginx:stable-alpine",
		ContainerPort: 80,
		HostPort:      9999,
		TestUrl:       "http://localhost:9999",
		Env: []string{
			"env1=value1",
			"evn2=value2",
		},
		Links: []string{
			"linkedContainer:link1",
		},
	}
	err = createContainer(cmp, logger.Infof)
	defer removeComponent(cmp, logger.Infof)
	if err != nil {
		t.Errorf("Error, expected container to be created, got %s", err.Error())
	}
}

func TestContainerWorkflow(t *testing.T) {
	logger := setUpForRun()
	if os.Getenv("NO_NETWORK") == "true" {
		t.Skipf("NO_NETWORK set to true, skipping")
	}
	common.SkipDockerTesting(t)
	cmp := common.Component{
		Name:     "test",
		DockerId: "test-container",
		Image:    "docker.io/library/nginx:stable-alpine",
	}

	err := createContainer(cmp, logger.Infof)
	defer removeComponent(cmp, logger.Infof)
	if err != nil {
		t.Errorf("Expected container to be created, got %s", err.Error())
	}

	err = stopContainer(cmp, logger.Infof)
	if err != nil {
		t.Errorf("Expected container to be stopped, got %s", err.Error())
	}

	err = startComponent(cmp, logger.Infof)
	if err != nil {
		t.Errorf("Expected container to be started, got %s", err.Error())
	}

	err = printLogs(cmp, false)
	if err != nil {
		t.Errorf("Expected container to print logs, got %s", err.Error())
	}

	err = removeComponent(cmp, logger.Infof)
	if err != nil {
		t.Errorf("Expected container to be removed, got %s", err.Error())
	}

	container, err := getContainer(cmp)
	if err == nil {
		t.Errorf("Expected container not to exist, got %s", container.Names)
	}
}

func TestDockerGetImages(t *testing.T) {
	common.SkipDockerTesting(t)
	images := getImages()
	if len(images) == 0 {
		t.Errorf("Expected to have at least some images")
	}
}

func Test_parseMounts(t *testing.T) {
	cwd, _ := os.Getwd()
	if _, err := parseMounts(common.Component{Mounts: []string{"invalid_format"}}); err == nil {
		t.Errorf("Expected error due to invalid format, got nothing")
	}
	if _, err := parseMounts(common.Component{Mounts: []string{"/non-existing:/"}}); err == nil {
		t.Errorf("Expected error due to non-existing source directory, got nothing")
	}
	mountsToParse := []string{filepath.Join("/") + ":/", "/etc:/etc"}
	if runtime.GOOS == "windows" {
		mountsToParse = []string{filepath.Join("/") + ":/", "/windows:/etc"}
	}

	mounts, err := parseMounts(common.Component{Mounts: mountsToParse})
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if len(mounts) != 2 {
		t.Errorf("Expected to get 2 mounts back, got %d", len(mounts))
	}
	expectedPath := filepath.VolumeName(cwd) + filepath.Join("/")
	if mounts[0].Source != expectedPath {
		t.Errorf("Expected first mount source to be '%s', got '%s'", expectedPath, mounts[0].Source)
	}
	if mounts[0].Target != "/" {
		t.Errorf("Expected first mount target to be '/', got '%s'", mounts[0].Target)
	}
}

func Test_mkContextTar(t *testing.T) {
	testRootDir, _ := setUpForBuild()
	defer os.RemoveAll(testRootDir)
	type args struct {
		contextDir string
		dockerFile string
	}
	tests := []struct {
		name         string
		args         args
		wantFileName bool
		wantErr      bool
	}{
		{
			name:         "test-pass",
			args:         args{contextDir: testRootDir + "/src", dockerFile: testRootDir + "/src/Dockerfile"},
			wantFileName: true,
			wantErr:      false,
		},
		{
			name:         "test-non-existing",
			args:         args{contextDir: "/non-existing"},
			wantFileName: false,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mkContextTar(tt.args.contextDir, tt.args.dockerFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("mkContextTar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantFileName && got == "" {
				t.Errorf("mkContextTar() = %v, wantFileName %v", got, tt.wantFileName)
			}
			if err == nil {
				// Check list of files in tar equals what's in the directory
				match := extractAndCompare(got, testRootDir)
				if !match {
					t.Errorf("Extracted tar contents don't match the source")
				}
			}
		})
	}
}

func Test_buildImage(t *testing.T) {
	mockDir, ctx := setUpForBuild()
	image := "test-image"
	buildRoot := mockDir + "/buildtest"
	dockerFile := mockDir + "/buildtest/Dockerfile"
	buildArgs := []string{"arg1:value1"}
	err := buildImage(ctx, image, buildRoot, dockerFile, buildArgs, true)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	// Test missing parameters
	if err := buildImage(ctx, "", "", "", buildArgs, true); err == nil {
		t.Errorf("Expected error (missing parameters), got nothing")
	}
	// Test invalid build root
	if err := buildImage(ctx, image, "../../../../../../../../../../../../../../", dockerFile, buildArgs, true); err == nil {
		t.Errorf("Expected error (invalid build root), got nothing")
	}
}

func Test_parseBuildArgs(t *testing.T) {
	os.Setenv("TEST_VAR", "value4")
	buildArgs := []string{"arg1:value1", "arg2:value2", "arg3:value3", "arg4:$TEST_VAR", "arg5", "arg6:", ":value7", ""}
	parsed := parseBuildArgs(buildArgs)
	if len(parsed) != 6 {
		t.Errorf("Unexpected length, expected 6, got %d", len(parsed))
	}
	if *parsed["arg1"] != "value1" {
		t.Errorf("Expected %s, got %s", "value1", *parsed["arg1"])
	}
	if *parsed["arg2"] != "value2" {
		t.Errorf("Expected %s, got %s", "value2", *parsed["arg2"])
	}
	if *parsed["arg3"] != "value3" {
		t.Errorf("Expected %s, got %s", "value3", *parsed["arg3"])
	}
	if *parsed["arg4"] != "value4" {
		t.Errorf("Expected %s, got %s", "value4", *parsed["arg4"])
	}
	if *parsed["arg5"] != "arg5" {
		t.Errorf("Expected %s, got %s", "arg5", *parsed["arg5"])
	}
	if *parsed["arg6"] != "" {
		t.Errorf("Expected %s, got %s", "", *parsed["arg6"])
	}
}

func extractAndCompare(tarFileName string, testRootDirectory string) bool {
	var cmd *exec.Cmd
	var untarOut, findSource, findDest []byte
	var err error

	// Extract tar to dest directory
	cmd = exec.Command("tar", "-xvf", tarFileName)
	cmd.Dir = testRootDirectory + "/dest"
	untarOut, err = cmd.Output()
	if err != nil {
		panic(err)
	}
	println(untarOut)

	cmd = exec.Command("find", ".")
	cmd.Dir = testRootDirectory + "/src"
	findSource, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command("find", ".")
	cmd.Dir = testRootDirectory + "/dest"
	findDest, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(testRootDirectory)
	fmt.Printf("Tar source directory: %s", findSource)
	fmt.Printf("Extracted tar contents: %s", findDest)

	return strings.Contains(string(findDest), string(findSource))
	//return reflect.DeepEqual(findSource, findDest)
}
