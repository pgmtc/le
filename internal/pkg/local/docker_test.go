package local

import (
	"testing"
)

<<<<<<< HEAD
func TestMissingParameters(t * testing.T) {
	cmp := Component {
		name: "test",
		dockerId: "test-container",
	}
	err := createContainer(cmp)
	if (err == nil) {
=======
func TestMissingParameters(t *testing.T) {
	cmp := Component{
		name:     "test",
		dockerId: "test-container",
	}
	err := createContainer(cmp)
	if err == nil {
>>>>>>> ebd35fcfdf40477b29a5c99a3629725738a8dfb3
		t.Errorf("Expected to fail due to mandatory missing")
	}
}

<<<<<<< HEAD

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
=======
func TestComplex(t *testing.T) {
	var err error

	cmp1 := Component{
		name:     "linkedContainer",
		dockerId: "linkedContainer",
		image:    "nginx:stable-alpine",
	}

	err = pullImage(cmp1)
	if err != nil {
		t.Errorf("Error, expected image to be pulled, got %s", err.Error())
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

	cmp := Component{
		name:          "test",
		dockerId:      "testContainer",
		image:         "bitnami/redis:latest",
		containerPort: 9999,
		hostPort:      9999,
		testUrl:       "http://localhost:8765/orchard-gateway-msvc/health",
		env: []string{
			"env1=value1",
			"evn2=value2",
		},
		links: []string{
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
	cmp := Component{
		name:     "test",
		dockerId: "test-container",
		image:    "nginx:stable-alpine",
	}

	err := createContainer(cmp)
	defer removeContainer(cmp)
	if err != nil {
>>>>>>> ebd35fcfdf40477b29a5c99a3629725738a8dfb3
		t.Errorf("Expected container to be created, got %s", err.Error())
	}

	err = stopContainer(cmp)
<<<<<<< HEAD
	if (err != nil) {
=======
	if err != nil {
>>>>>>> ebd35fcfdf40477b29a5c99a3629725738a8dfb3
		t.Errorf("Expected container to be stopped, got %s", err.Error())
	}

	err = startContainer(cmp)
<<<<<<< HEAD
	if (err != nil) {
		t.Errorf("Expected container to be started, got %s", err.Error())
	}

	err = removeContainer(cmp)
	if  (err != nil) {
=======
	if err != nil {
		t.Errorf("Expected container to be started, got %s", err.Error())
	}

	err = dockerPrintLogs(cmp, false)
	if err != nil {
		t.Errorf("Expected container to print logs, got %s", err.Error())
	}

	err = removeContainer(cmp)
	if err != nil {
>>>>>>> ebd35fcfdf40477b29a5c99a3629725738a8dfb3
		t.Errorf("Expected container to be removed, got %s", err.Error())
	}

	container, err := getContainer(cmp)
<<<<<<< HEAD
	if (err == nil) {
		t.Errorf("Expected container not to exist, got %s", container.Names)
	}
}
=======
	if err == nil {
		t.Errorf("Expected container not to exist, got %s", container.Names)
	}
}

func Test_pullImage(t *testing.T) {
	type args struct {
		component Component
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "testNginx",
			args: args{
				component: Component{
					name:  "ironGo",
					image: "iron/go",
				},
			},
			wantErr: false,
		},
		{
			name: "testNonExisting",
			args: args{
				component: Component{
					name:  "nonExisting",
					image: "whatever-nonexisting",
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
>>>>>>> ebd35fcfdf40477b29a5c99a3629725738a8dfb3
