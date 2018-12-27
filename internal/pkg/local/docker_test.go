package local

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"testing"
)

func TestMissingParameters(t *testing.T) {
	cmp := common.Component{
		Name:     "test",
		DockerId: "test-container",
	}
	err := createContainer(cmp)
	if err == nil {
		t.Errorf("Expected to fail due to mandatory missing")
	}
}

func TestComplex(t *testing.T) {
	var err error

	cmp1 := common.Component{
		Name:     "linkedContainer",
		DockerId: "linkedContainer",
		Image:    "nginx:stable-alpine",
	}

	err = pullImage(cmp1)
	if err != nil {
		t.Errorf("Error, expected Image to be pulled, got %s", err.Error())
	}

	err = createContainer(cmp1)
	defer removeContainer(cmp1)
	if err != nil {
		t.Errorf("Error, expected container to be created, got %s", err.Error())
	}
	err = startContainer(cmp1)
	if err != nil {
		t.Errorf("Error, expected container to be created, got %s", err.Error())
	}

	cmp := common.Component{
		Name:          "test",
		DockerId:      "testContainer",
		Image:         "bitnami/redis:latest",
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
	err = createContainer(cmp)
	defer removeContainer(cmp)
	if err != nil {
		t.Errorf("Error, expected container to be created, got %s", err.Error())
	}
}

func TestContainerWorkflow(t *testing.T) {
	cmp := common.Component{
		Name:     "test",
		DockerId: "test-container",
		Image:    "nginx:stable-alpine",
	}

	err := createContainer(cmp)
	defer removeContainer(cmp)
	if err != nil {
		t.Errorf("Expected container to be created, got %s", err.Error())
	}

	err = stopContainer(cmp)
	if err != nil {
		t.Errorf("Expected container to be stopped, got %s", err.Error())
	}

	err = startContainer(cmp)
	if err != nil {
		t.Errorf("Expected container to be started, got %s", err.Error())
	}

	err = dockerPrintLogs(cmp, false)
	if err != nil {
		t.Errorf("Expected container to print logs, got %s", err.Error())
	}

	err = removeContainer(cmp)
	if err != nil {
		t.Errorf("Expected container to be removed, got %s", err.Error())
	}

	container, err := getContainer(cmp)
	if err == nil {
		t.Errorf("Expected container not to exist, got %s", container.Names)
	}
}

func Test_pullImage(t *testing.T) {
	type args struct {
		component common.Component
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "testNginx",
			args: args{
				component: common.Component{
					Name:  "ironGo",
					Image: "iron/go",
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
			if err := pullImage(tt.args.component); (err != nil) != tt.wantErr {
				t.Errorf("pullImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
