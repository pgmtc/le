package docker

import (
	"github.com/pgmtc/le/pkg/common"
)

type Builder struct{}

func (Builder) BuildImage(ctx common.Context, image string, buildRoot string, dockerFile string, buildArgs []string, noCache bool) error {
	return buildImage(ctx, image, buildRoot, dockerFile, buildArgs, noCache)
}
