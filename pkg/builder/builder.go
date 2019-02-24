package builder

import (
	"github.com/pgmtc/le/pkg/common"
	"github.com/pkg/errors"
)

type Builder interface {
	BuildImage(ctx common.Context, image string, buildRoot string, dockerFile string, buildArgs []string, noCache bool) error
}

type SpyBuilder struct {
	WantErrorMessage string
	Spy              struct {
		Ctx        common.Context
		Image      string
		BuildRoot  string
		DockerFile string
		BuildArgs  []string
		NoCache    bool
	}
}

func (m *SpyBuilder) BuildImage(ctx common.Context, image string, buildRoot string, dockerFile string, buildArgs []string, noCache bool) error {
	m.Spy.Image = image
	m.Spy.BuildRoot = buildRoot
	m.Spy.DockerFile = dockerFile
	m.Spy.BuildArgs = buildArgs
	m.Spy.NoCache = noCache
	m.Spy.Ctx = ctx

	if m.WantErrorMessage != "" {
		return errors.Errorf(m.WantErrorMessage)
	}
	return nil
}
