package common

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
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
	usr, _ := user.Current()
	path = strings.Replace(path, "~", usr.HomeDir, 1)

	result, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return
}

func YamlMarshall(data interface{}, fileName string) (resultErr error) {
	bytes, _ := yaml.Marshal(data)
	if err := ioutil.WriteFile(fileName, bytes, 0644); err != nil {
		resultErr = errors.Errorf("error writing file: %s", err.Error())
		return
	}

	return
}

func YamlUnmarshall(fileName string, out interface{}) (resultErr error) {
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
