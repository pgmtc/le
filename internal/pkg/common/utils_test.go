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

func Test_YamlMarshall(t *testing.T) {
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

	tmpFile := tmpDir + "/YamlMarshall-out.yaml"
	junkFile := tmpDir + "/junk.yaml"
	ioutil.WriteFile(junkFile, []byte("junk file content\n"), 0644)

	// Test write success
	err := YamlMarshall(testData, tmpFile)
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	// Test write failure (non existing directory)
	err = YamlMarshall(testData, "/non-existing"+tmpFile)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

	outData := TestType{}
	err = YamlUnmarshall(tmpFile, &outData)
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	if !reflect.DeepEqual(testData, outData) {
		t.Errorf("In and out does not match")
	}

	// Test Junk Content
	err = YamlUnmarshall(junkFile, &outData)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Test Non Existing file
	err = YamlUnmarshall(junkFile+"-non-existing", &outData)
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
