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
	Run(log Logger, args ...string) error
}
