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
}

func TestDockerRunner_Pull(t *testing.T) {
	if err := runner.Pull(ctx, ctx.Config.CurrentProfile().Components[0]); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestDockerRunner_Create(t *testing.T) {
	if err := runner.Create(ctx, ctx.Config.CurrentProfile().Components[0]); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestDockerRunner_Start(t *testing.T) {
	if err := runner.Start(ctx, ctx.Config.CurrentProfile().Components[0]); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestDockerRunner_Stop(t *testing.T) {
	if err := runner.Stop(ctx, ctx.Config.CurrentProfile().Components[0]); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestDockerRunner_Logs(t *testing.T) {
	if err := runner.Logs(ctx, ctx.Config.CurrentProfile().Components[0]); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestDockerRunner_Remove(t *testing.T) {
	if err := runner.Remove(ctx, ctx.Config.CurrentProfile().Components[0]); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}
