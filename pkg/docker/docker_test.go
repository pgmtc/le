package docker

import (
	"os"
	"strings"
	"testing"

	"github.com/pgmtc/le/pkg/common"
)

func setUp() (logger common.Logger) {
	logger = &common.StringLogger{}
	return
}

func TestMissingParameters(t *testing.T) {
	logger := setUp()
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
	logger := setUp()
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
				images := dockerGetImages()
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
				images = dockerGetImages()
				if common.ArrContains(images, tt.args.component.Image) {
					t.Errorf("Pulled image '%s' still exist, should have been removed", tt.args.component.Image)
				}
			}

		})
	}
}

func TestComplex(t *testing.T) {
	logger := setUp()
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
	logger := setUp()
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

	err = dockerPrintLogs(cmp, false)
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
	images := dockerGetImages()
	if len(images) == 0 {
		t.Errorf("Expected to have at least some images")
	}
}
