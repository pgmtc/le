package local

import (
	"testing"

	"github.com/pgmtc/orchard-cli/internal/pkg/common"
)

func TestMissingParameters(t *testing.T) {
	cmp := common.Component{
		Name:     "test",
		DockerId: "test-container",
	}
	err := createContainer(cmp, common.HandlerArguments{})
	if err == nil {
		t.Errorf("Expected to fail due to mandatory missing")
	}
}

func Test_pullImage(t *testing.T) {
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
			name: "testIronGo",
			args: args{
				component: common.Component{
					Name:  "ironGo",
					Image: "iron/go:latest",
				},
			},
			wantErr: false,
		},
		{
			name: "testNonExisting",
			args: args{
				component: common.Component{
					Name:  "nonExisting",
					Image: "whatever-nonexisting",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			removeImage(tt.args.component, common.HandlerArguments{}) // Ignore errors, image may not exist

			err := pullImage(tt.args.component, common.HandlerArguments{})
			if (err != nil) != tt.wantErr {
				t.Errorf("pullImage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				// Check that the image exists
				images, err := dockerGetImages()
				if err != nil {
					t.Errorf("Unexpected error when getting list of images: %s", err.Error())
				}

				if !common.ArrContains(images, tt.args.component.Image) {
					t.Errorf("Pulled image '%s' seems not to exist", tt.args.component.Image)
				}

				// Try to remove image
				err = removeImage(tt.args.component, common.HandlerArguments{})
				if err != nil {
					t.Errorf("Unexpected error when removing image: %s", err.Error())
				}
				images, err = dockerGetImages()
				if err != nil {
					t.Errorf("Unexpected error when getting list of images: %s", err.Error())
				}
				if common.ArrContains(images, tt.args.component.Image) {
					t.Errorf("Pulled image '%s' still exist, should have been removed", tt.args.component.Image)
				}
			}

		})
	}
}

func TestComplex(t *testing.T) {
	common.SkipDockerTesting(t)
	var err error

	cmp1 := common.Component{
		Name:     "linkedContainer",
		DockerId: "linkedContainer",
		Image:    "nginx:stable-alpine",
	}

	err = pullImage(cmp1, common.HandlerArguments{})
	if err != nil {
		t.Errorf("Error, expected Image to be pulled, got %s", err.Error())
	}

	err = createContainer(cmp1, common.HandlerArguments{})
	defer removeContainer(cmp1, common.HandlerArguments{})
	if err != nil {
		t.Errorf("Error, expected container to be created, got %s", err.Error())
	}
	err = startContainer(cmp1, common.HandlerArguments{})
	if err != nil {
		t.Errorf("Error, expected container to be created, got %s", err.Error())
	}

	cmp := common.Component{
		Name:          "test",
		DockerId:      "testContainer",
		Image:         "nginx:stable-alpine",
		ContainerPort: 9999,
		HostPort:      9999,
		TestUrl:       "http://localhost:8765/orchard-gateway-msvc/health",
		Env: []string{
			"env1=value1",
			"evn2=value2",
		},
		Links: []string{
			"linkedContainer:link1",
		},
	}
	err = createContainer(cmp, common.HandlerArguments{})
	defer removeContainer(cmp, common.HandlerArguments{})
	if err != nil {
		t.Errorf("Error, expected container to be created, got %s", err.Error())
	}
}

func TestContainerWorkflow(t *testing.T) {
	common.SkipDockerTesting(t)
	cmp := common.Component{
		Name:     "test",
		DockerId: "test-container",
		Image:    "nginx:stable-alpine",
	}

	err := createContainer(cmp, common.HandlerArguments{})
	defer removeContainer(cmp, common.HandlerArguments{})
	if err != nil {
		t.Errorf("Expected container to be created, got %s", err.Error())
	}

	err = stopContainer(cmp, common.HandlerArguments{})
	if err != nil {
		t.Errorf("Expected container to be stopped, got %s", err.Error())
	}

	err = startContainer(cmp, common.HandlerArguments{})
	if err != nil {
		t.Errorf("Expected container to be started, got %s", err.Error())
	}

	err = dockerPrintLogs(cmp, false)
	if err != nil {
		t.Errorf("Expected container to print logs, got %s", err.Error())
	}

	err = removeContainer(cmp, common.HandlerArguments{})
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
	if _, err := dockerGetImages(); err != nil {
		t.Errorf("Unexpected error, but got %s", err.Error())
	}
}
