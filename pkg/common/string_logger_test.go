package common

import "testing"

var log = StringLogger{}

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

func TestStringLogger_Write(t *testing.T) {
	preCount := len(log.InfoMessages)
	log.Write([]byte("test-message"))
	postCount := len(log.InfoMessages)

	if preCount != (postCount - 1) {
		t.Errorf("Expected that there will be an extra info message after write()")
	}
	lastMessage := log.InfoMessages[len(log.InfoMessages)-1]
	if lastMessage != "test-message" {
		t.Errorf("Expected last message to be %s, got %s", "test-message", log.InfoMessages[len(log.InfoMessages)-1])
	}
}
