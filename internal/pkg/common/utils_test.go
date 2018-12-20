package common

import (
	"testing"
)

func TestArrContains(t *testing.T) {
	arr := []string {"element1", "element2", "element3"}
	var emptyArr []string

	if !ArrContains(arr, "element1") {
		t.Errorf("Expected true to be returned")
	}
	if ArrContains(arr, "nonExistent") {
		t.Errorf("Expected false to be returned")
	}
	if ArrContains(arr, "") {
		t.Errorf("Expected false to be returned")
	}
	if ArrContains(emptyArr, "element2") {
		t.Errorf("Expected false to be returned")
	}

}
