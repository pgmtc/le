package source

import "testing"

func Test_pull(t *testing.T) {
	err := pull([]string {})
	if (err != nil) {
		t.Errorf("Unexpected error returned: %s", err.Error())
	}
}
