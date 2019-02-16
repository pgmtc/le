package common

import "testing"

var log = StringLogger{}
var msg1 = "first message: %s"
var msg2 = "second message: %s"
var msg3 = "third message: %s"
var msg4 = "fourth message: %s"

func testLogMethod(t *testing.T, logMethod func(format string, a ...interface{}), messages *[]string) {
	values := [][]string{
		{"first message: %s", "one", "first message: one"},
		{"second message: %s", "two", "second message: two"},
		{"third message: %s", "three", "third message: three"},
		{"fourth message: %s", "four", "fourth message: four"},
	}
	for i, row := range values {
		logMethod(row[0], row[1])
		if len(*messages) != i+1 {
			t.Errorf("Expected to have %d messages in total. Have %d", i+1, len(log.ErrorMessages))
		}
		if (*messages)[i] != row[2] {
			t.Errorf("Expected message 1 to be %s, got %s", row[2], log.ErrorMessages[i])
		}
	}
}

func TestStringLogger_Errorf(t *testing.T) {
	testLogMethod(t, log.Errorf, &log.ErrorMessages)
}

func TestStringLogger_Debugf(t *testing.T) {
	testLogMethod(t, log.Debugf, &log.DebugMessages)
}

func TestStringLogger_Infof(t *testing.T) {
	testLogMethod(t, log.Infof, &log.InfoMessages)
}
