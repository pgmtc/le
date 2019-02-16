package common

import (
	"fmt"
	"github.com/fatih/color"
)

type StringLogger struct {
	DebugMessages []string
	InfoMessages  []string
	ErrorMessages []string
}

func (l *StringLogger) Errorf(format string, a ...interface{}) {
	l.ErrorMessages = append(l.ErrorMessages, fmt.Sprintf(color.HiRedString(format, a...)))
}

func (l *StringLogger) Debugf(format string, a ...interface{}) {
	l.DebugMessages = append(l.DebugMessages, fmt.Sprintf(color.HiRedString(format, a...)))
}

func (l *StringLogger) Infof(format string, a ...interface{}) {
	l.InfoMessages = append(l.InfoMessages, fmt.Sprintf(color.HiRedString(format, a...)))
}
