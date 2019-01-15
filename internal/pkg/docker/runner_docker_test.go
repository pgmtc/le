package docker

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"testing"
)

var runner = DockerRunner{}
var ctx = common.Context{
	Config: common.CreateMockConfig([]common.Component{
		{
			Name:     "test-component",
			DockerId: "test-component",
			Image:    "iron/go",
		},
	}),
	Log: common.ConsoleLogger{},
}

func TestDockerRunner_Status(t *testing.T) {
	runner.Status(ctx)
	runner.Status(ctx, "-v")
}

func TestDockerRunner_Pull(t *testing.T) {
	cmp := ctx.Config.CurrentProfile().Components[0]
	if err := runner.Pull(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	defer removeImage(cmp)
}

func TestDockerRunner_Create(t *testing.T) {
	cmp := ctx.Config.CurrentProfile().Components[0]
	runner.Remove(ctx, cmp)
	if err := runner.Create(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	defer runner.Remove(ctx, cmp)
}

func TestDockerRunner_Start(t *testing.T) {
	cmp := ctx.Config.CurrentProfile().Components[0]
	if err := runner.Start(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	defer runner.Stop(ctx, cmp)
}

func TestDockerRunner_Stop(t *testing.T) {
	cmp := ctx.Config.CurrentProfile().Components[0]
	runner.Start(ctx, cmp)
	if err := runner.Stop(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestDockerRunner_Logs(t *testing.T) {
	cmp := ctx.Config.CurrentProfile().Components[0]
	if err := runner.Logs(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestDockerRunner_Remove(t *testing.T) {
	cmp := ctx.Config.CurrentProfile().Components[0]
	if err := runner.Remove(ctx, cmp); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}
