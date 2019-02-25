package local

import "github.com/pgmtc/le/pkg/common"

//go:generate mockgen -destination=./mocks/mock_runner.go -package=mocks github.com/pgmtc/le/pkg/local Runner
type Runner interface {
	Create(ctx common.Context, cmp common.Component) error
	Remove(ctx common.Context, cmp common.Component) error
	Start(ctx common.Context, cmp common.Component) error
	Stop(ctx common.Context, cmp common.Component) error
	Pull(ctx common.Context, cmp common.Component) error
	Logs(ctx common.Context, cmp common.Component, follow bool) error
	Status(ctx common.Context, args ...string) error
}
