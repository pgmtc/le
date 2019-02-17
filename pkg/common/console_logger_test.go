package common

import "testing"

// TODO: Implement more advanced tests
func TestConsoleLogger_Errorf(t *testing.T) {
	c := ConsoleLogger{}
	c.Errorf("%s", "Error message")
}

func TestConsoleLogger_Debugf(t *testing.T) {
	c := ConsoleLogger{}
	c.Debugf("%s", "Debug message")
}

func TestConsoleLogger_Infof(t *testing.T) {
	c := ConsoleLogger{}
	c.Infof("%s", "Info message")
}

func TestConsoleLogger_Write(t *testing.T) {
	c := ConsoleLogger{}
	c.Write([]byte("Test message"))
}
