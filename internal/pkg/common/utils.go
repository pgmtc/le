package common

import (
	"os"
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
