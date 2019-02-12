package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

func mockContext() common.Context {
	config := common.CreateMockConfig([]common.Component{})
	ctx := common.Context{
		Config: config,
		Log:    &common.StringLogger{},
	}
	return ctx
}

func mockupDir() string {
	tmpDir, _ := ioutil.TempDir("", "orchard-test-mock")
	os.MkdirAll(tmpDir+"/src", os.ModePerm)
	os.MkdirAll(tmpDir+"/dest", os.ModePerm)

	os.MkdirAll(tmpDir+"/src/subdir", os.ModePerm)
	os.MkdirAll(tmpDir+"/src/.hiddendir", os.ModePerm)

	fileContent := []byte("testing file contents\n")
	ioutil.WriteFile(tmpDir+"/src/"+"file1.txt", fileContent, 0644)
	ioutil.WriteFile(tmpDir+"/src/subdir/"+"file2.txt", fileContent, 0644)
	ioutil.WriteFile(tmpDir+"/src/.hiddendir/"+"file3.txt", fileContent, 0644)
	ioutil.WriteFile(tmpDir+"/src/Dockerfile", []byte("FROM scratch"), 0644)

	os.MkdirAll(tmpDir+"/jartest1/target", os.ModePerm)
	ioutil.WriteFile(tmpDir+"/jartest1/target/test-file-1.jar", fileContent, 0644)

	os.MkdirAll(tmpDir+"/jartest2/target", os.ModePerm)
	ioutil.WriteFile(tmpDir+"/jartest2/target/test-file-1.jar", fileContent, 0644)
	ioutil.WriteFile(tmpDir+"/jartest2/target/test-file-2.jar", fileContent, 0644)

	os.MkdirAll(tmpDir+"/jartest3/target", os.ModePerm)
	ioutil.WriteFile(tmpDir+"/jartest3/target/test-file-1.notjar", fileContent, 0644)

	os.MkdirAll(tmpDir+"/jartest4", os.ModePerm)
	ioutil.WriteFile(tmpDir+"/jartest4/test-file-1.jar", fileContent, 0644) // not in target subdirectory

	os.MkdirAll(tmpDir+"/buildtest", os.ModePerm)
	ioutil.WriteFile(tmpDir+"/buildtest/Dockerfile", []byte("FROM scratch\nADD . ."), 0644)
	ioutil.WriteFile(tmpDir+"/buildtest/Dockerfile-invalid", []byte("FROM rubbish\nADD . ."), 0644)
	ioutil.WriteFile(tmpDir+"/buildtest/config.yaml", []byte(""+
		"image: test-image\n"+
		"buildroot: "+tmpDir+"/buildtest/\n"+
		"dockerfile: "+tmpDir+"/buildtest/Dockerfile\n"), 0644)

	return tmpDir
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

	return reflect.DeepEqual(findSource, findDest)
}

func Test_build(t *testing.T) {
	common.SkipDockerTesting(t)
	//testingRoot := mockupDir()
	t.Skip("WIP")
}

func Test_mkContextTar(t *testing.T) {
	testRootDir := mockupDir()
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

func Test_relativeOrAbsolute(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name        string
		args        args
		wantChanged bool
	}{
		{
			name:        "test-relative",
			args:        args{path: "relative/path.txt"},
			wantChanged: true,
		},
		{
			name:        "test-absolute",
			args:        args{path: "/absolute/relative/path.txt"},
			wantChanged: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := common.ParsePath(tt.args.path)
			if tt.wantChanged && tt.args.path == got {
				t.Errorf("relativeOrAbsolute() = %v, want %v", got, tt.args.path)
			}
			if !tt.wantChanged && tt.args.path != got {
				t.Errorf("relativeOrAbsolute() = %v, want %v", got, tt.args.path)
			}
		})
	}
}

func Test_findMsvcJar(t *testing.T) {
	testRootDir := mockupDir()
	defer os.RemoveAll(testRootDir)
	type args struct {
		path string
	}
	tests := []struct {
		name         string
		args         args
		wantFileName string
		wantErr      bool
	}{
		{
			name:         "test-one-exist",
			args:         args{path: testRootDir + "/jartest1"},
			wantFileName: "target/test-file-1.jar",
			wantErr:      false,
		},
		{
			name:         "test-two-exist",
			args:         args{path: testRootDir + "/jartest2"},
			wantFileName: "",
			wantErr:      true,
		},
		{
			name:         "test-none-exist",
			args:         args{path: testRootDir + "/jartest3"},
			wantFileName: "",
			wantErr:      false,
		},
		{
			name:         "test-no-target-directory",
			args:         args{path: testRootDir + "/jartest4"},
			wantFileName: "",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileName, err := findMsvcJar(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("findMsvcJar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFileName != tt.wantFileName {
				t.Errorf("findMsvcJar() = %v, want %v", gotFileName, tt.wantFileName)
			}
		})
	}
}

func Test_parseBuildProperties(t *testing.T) {
	tmpDir := mockupDir()
	buildDir := tmpDir + "/buildtest"
	err, image, buildDir, dockerFile := parseBuildProperties(buildDir)
	expectedImage := "some-image"
	expectedBuildDir := "some-builddir"
	expectedDockerfile := "some-dockerfile"
	if err != nil {
		t.Errorf("Unexpected error returned: %s", err.Error())
	}
	if image == "" {
		t.Errorf("Expected %s, got %s", expectedImage, image)
	}
	if buildDir == "" {
		t.Errorf("Expected %s, got %s", expectedBuildDir, buildDir)
	}
	if dockerFile == "" {
		t.Errorf("Expected %s, got %s", expectedDockerfile, dockerFile)
	}

	// Test error
	buildDir = tmpDir + "/non-existing"
	err, image, buildDir, dockerFile = parseBuildProperties(buildDir)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

}

func Test_buildImage(t *testing.T) {
	mockDir := mockupDir()
	ctx := mockContext()
	image := "test-image"
	buildRoot := mockDir + "/buildtest"
	dockerFile := mockDir + "/buildtest/Dockerfile"
	noCache := true
	err := buildImage(ctx, image, buildRoot, dockerFile, noCache)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func Test_buildAction(t *testing.T) {
	buildContext := mockupDir() + "/buildtest"
	if err := buildAction.Run(mockContext(), "--specdir", buildContext, "--nocache"); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	// Try missing specdir
	if err := buildAction.Run(mockContext(), "--specdir"); err == nil {
		t.Errorf("Expected error, got nothing")
	}

}
