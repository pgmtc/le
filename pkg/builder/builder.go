package builder

import (
	"github.com/pgmtc/le/pkg/common"
)

//go:generate mockgen -destination=./mocks/mock_builder.go -package=mocks github.com/pgmtc/le/pkg/builder Builder
type Builder interface {
	BuildImage(ctx common.Context, image string, buildRoot string, dockerFile string, buildArgs []string, noCache bool) error
}
