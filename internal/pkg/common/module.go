package common

type Module interface {
	//Run(Log Logger, args ...string) error
	GetActions() map[string]Action
}

type Logger interface {
	Errorf(format string, a ...interface{})
	Debugf(format string, a ...interface{})
	Infof(format string, a ...interface{})
}

type Context struct {
	Log    Logger
	Config Configuration
	Module Module
}

type Action interface {
	Run(ctx Context, args ...string) error
}

type RawAction struct {
	Handler func(ctx Context, args ...string) error
}

func (a *RawAction) Run(ctx Context, args ...string) error {
	return a.Handler(ctx, args...)
}
