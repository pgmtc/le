package local

import (
	"testing"
)

func TestMissingParameters(t * testing.T) {
	cmp := Component {
		name: "test",
		dockerId: "test-container",
	}
	err := createContainer(cmp)
	if (err == nil) {
		t.Errorf("Expected to fail due to mandatory missing")
	}
}


func TestComplex(t *testing.T) {
	cmp := Component {
		name: "test",
		dockerId: "test-container",
		image: "bitnami/redis:latest",
		containerPort: 8080,
		hostPort: 8765,
		testUrl:  "http://localhost:8765/orchard-gateway-msvc/health",
	}
	err := createContainer(cmp)
	if (err != nil) {
		t.Errorf("Expected to fail due to mandatory missing")
	}

	removeContainer(cmp)
}

func TestSimpleContainerWorkflow(t *testing.T) {
	cmp := Component {
		name: "test",
		dockerId: "test-container",
		image: "bitnami/redis:latest",
	}

	err := createContainer(cmp)
	if (err != nil) {
		t.Errorf("Expected container to be created, got %s", err.Error())
	}

	err = stopContainer(cmp)
	if (err != nil) {
		t.Errorf("Expected container to be stopped, got %s", err.Error())
	}

	err = startContainer(cmp)
	if (err != nil) {
		t.Errorf("Expected container to be started, got %s", err.Error())
	}

	err = removeContainer(cmp)
	if  (err != nil) {
		t.Errorf("Expected container to be removed, got %s", err.Error())
	}

	container, err := getContainer(cmp)
	if (err == nil) {
		t.Errorf("Expected container not to exist, got %s", container.Names)
	}
}