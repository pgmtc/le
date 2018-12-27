package builder

import "testing"

func Test_build(t *testing.T) {
	err := build([]string{})
	if err != nil {
		t.Errorf("Unexpected error returned: %s", err.Error())
	}
}
