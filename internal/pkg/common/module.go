package common

type Module interface {
	//Run(log Logger, args ...string) error
	GetActions() map[string]Action
}

type Logger interface {
	Errorf(format string, a ...interface{})
	Debugf(format string, a ...interface{})
	Infof(format string, a ...interface{})
}

type Action interface {
	Run(log Logger, config Configuration, args ...string) error
}

type RawAction struct {
	Handler func(log Logger, config Configuration, args ...string) error
}

func (a *RawAction) Run(log Logger, config Configuration, args ...string) error {
	return a.Handler(log, config, args...)
}
