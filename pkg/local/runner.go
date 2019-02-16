package local

import "github.com/pgmtc/le/pkg/common"

type Runner interface {
	Create(ctx common.Context, cmp common.Component) error
	Remove(ctx common.Context, cmp common.Component) error
	Start(ctx common.Context, cmp common.Component) error
	Stop(ctx common.Context, cmp common.Component) error
	Pull(ctx common.Context, cmp common.Component) error
	Logs(ctx common.Context, cmp common.Component, follow bool) error
	Status(ctx common.Context, args ...string) error
}
