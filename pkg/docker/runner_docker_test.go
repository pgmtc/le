package docker

import (
	"github.com/pgmtc/le/pkg/common"
	"os"
	"testing"
)

var runner = Runner{}
var ctx = common.Context{
	Config: common.CreateMockConfig([]common.Component{
		{
			Name:          "test-component",
			DockerId:      "test-component",
			Image:         "docker.io/library/nginx:stable-alpine",
			ContainerPort: 80,
			HostPort:      9998,
			TestUrl:       "http://localhost:9998",
		},
		{
			Name:     "test-internal-component",
			DockerId: "test-component",
			Image:    "674155361995.dkr.ecr.eu-west-1.amazonaws.com/orchard/orchard-valuation-client-ui:latest",
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
	logger := setUp()
	if os.Getenv("NO_NETWORK") == "true" {
		t.Skipf("NO_NETWORK set to true, skipping")
	}
	cmp := ctx.Config.CurrentProfile().Components[0]
	if err := runner.Pull(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	defer removeImage(cmp, logger.Infof)
}

func TestDockerRunner_Workflow(t *testing.T) {
	logger := setUp()
	if os.Getenv("NO_NETWORK") == "true" {
		t.Skipf("NO_NETWORK set to true, skipping")
	}
	cmp := ctx.Config.CurrentProfile().Components[0]
	internalCmp := ctx.Config.CurrentProfile().Components[1]
	invalidCmp := ctx.Config.CurrentProfile().Components[2]

	runner.Remove(ctx, cmp)        // Don't handle errors - prepare for the test
	removeImage(cmp, logger.Infof) // Don't handle errors - prepare for the test

	if err := runner.Pull(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	defer removeImage(cmp, logger.Infof)

	runner.Remove(ctx, internalCmp)        // Don't handle errors - prepare for the test
	removeImage(internalCmp, logger.Infof) // Don't handle errors - prepare for the test
	if err := runner.Pull(ctx, internalCmp); err != nil && (os.Getenv("SKIP_AWS_TESTING") == "") {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	defer removeImage(internalCmp, logger.Infof)

	runner.Remove(ctx, invalidCmp)        // Don't handle errors - prepare for the test
	removeImage(invalidCmp, logger.Infof) // Don't handle errors - prepare for the test
	if err := runner.Pull(ctx, invalidCmp); err == nil {
		t.Errorf("Expected error, got nothing")
	}
	defer removeImage(invalidCmp, logger.Infof)

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
