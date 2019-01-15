package docker

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"os"
	"testing"
)

var runner = DockerRunner{}
var ctx = common.Context{
	Config: common.CreateMockConfig([]common.Component{
		{
			Name:          "test-component",
			DockerId:      "test-component",
			Image:         "nginx:stable-alpine",
			ContainerPort: 80,
			HostPort:      9999,
			TestUrl:       "http://localhost:9999",
		},
		{
			Name:     "test-internal-component",
			DockerId: "test-component",
			Image:    "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-local-db:latest",
		},
		{
			Name:     "test-invalid",
			DockerId: "test-component",
		},
	}),
	Log: common.ConsoleLogger{},
}

func TestDockerRunner_Status(t *testing.T) {
	runner.Status(ctx)
	runner.Status(ctx, "-v")
	runner.Status(ctx, "-f", "2")
	runner.Status(ctx, "-v", "-f", "2")
}

func TestDockerRunner_Pull(t *testing.T) {
	cmp := ctx.Config.CurrentProfile().Components[0]
	if err := runner.Pull(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	defer removeImage(cmp)
}

func TestDockerRunner_Workflow(t *testing.T) {
	cmp := ctx.Config.CurrentProfile().Components[0]
	internalCmp := ctx.Config.CurrentProfile().Components[1]
	invalidCmp := ctx.Config.CurrentProfile().Components[2]

	runner.Remove(ctx, cmp) // Don't handle errors - prepare for the test
	removeImage(cmp)        // Don't handle errors - prepare for the test

	if err := runner.Pull(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	defer removeImage(cmp)

	runner.Remove(ctx, internalCmp) // Don't handle errors - prepare for the test
	removeImage(internalCmp)        // Don't handle errors - prepare for the test
	if err := runner.Pull(ctx, internalCmp); err != nil && (os.Getenv("SKIP_AWS_TESTING") == "") {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	defer removeImage(internalCmp)

	runner.Remove(ctx, invalidCmp) // Don't handle errors - prepare for the test
	removeImage(invalidCmp)        // Don't handle errors - prepare for the test
	if err := runner.Pull(ctx, invalidCmp); err == nil {
		t.Errorf("Expected error, got nothing")
	}
	defer removeImage(invalidCmp)

	if err := runner.Create(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	// Recreate again - should get error
	if err := runner.Create(ctx, cmp); err == nil {
		t.Errorf("Expected error, got nothing")
	}
	defer runner.Remove(ctx, cmp)

	if err := runner.Start(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if err := runner.Logs(ctx, cmp, false); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if err := runner.Stop(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	runner.Create(ctx, cmp) // Don't handle errors
	if err := runner.Remove(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

}
