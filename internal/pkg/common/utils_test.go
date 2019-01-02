package common

import (
	"os"
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

func TestSkipDockerTesting_false(t *testing.T) {
	dockerTestingHappened = false
	origValue := os.Getenv("SKIP_DOCKER_TESTING")
	os.Unsetenv("SKIP_DOCKER_TESTING")
	defer evalSkipDockerTesting(t, true, origValue)
	SkipDockerTesting(t)
	dockerTestingHappened = true // It should get here
}

func evalSkipDockerTesting(t *testing.T, expectedValue bool, origValue string) {
	if (origValue == "") {
		os.Unsetenv("SKIP_DOCKER_TESTING")
	} else {
		os.Setenv("SKIP_DOCKER_TESTING", origValue)
	}

	if expectedValue != dockerTestingHappened {
		t.Errorf("Unexpected dockerTestingHappened value, expected %t got %t", expectedValue, dockerTestingHappened)
	}
}
