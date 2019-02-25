package common

type Module interface {
	//Run(Log Logger, args ...string) error
	GetActions() map[string]Action
}

type Logger interface {
	Errorf(format string, a ...interface{})
	Debugf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Write(p []byte) (n int, err error)
}

type Context struct {
	Log    Logger
	Config ConfigProvider
	Module Module
}

type Action interface {
	Run(ctx Context, args ...string) error
}

type RawActionhandler func(ctx Context, args ...string) error

type RawAction struct {
	Handler RawActionhandler
}

func (a *RawAction) Run(ctx Context, args ...string) error {
	return a.Handler(ctx, args...)
}
