package common

import (
	"fmt"
	"github.com/fatih/color"
)

type ConsoleLogger struct{}

func (ConsoleLogger) Errorf(format string, a ...interface{}) {
	fmt.Printf(color.HiRedString(format, a...))
}

func (ConsoleLogger) Debugf(format string, a ...interface{}) {
	fmt.Printf(color.HiBlackString(format, a...))
}

func (ConsoleLogger) Infof(format string, a ...interface{}) {
	fmt.Printf(color.WhiteString(format, a...))
}
