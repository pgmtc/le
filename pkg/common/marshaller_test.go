package common

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestYamlMarshaller_Marshall(t *testing.T) {
	marshaller := YamlMarshaller{}
	type TestType struct {
		Name  string
		Value int
	}

	testData := TestType{
		Name:  "Test",
		Value: 10,
	}

	tmpDir, _ := ioutil.TempDir("", "le-test-mock")
	defer os.RemoveAll(tmpDir)

	tmpFile := tmpDir + "/YamlMarshall-out.yaml"
	junkFile := tmpDir + "/junk.yaml"
	ioutil.WriteFile(junkFile, []byte("junk file content\n"), 0644)

	// Test write success
	err := marshaller.Marshall(testData, tmpFile)
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	// Test write failure (non existing directory)
	err = marshaller.Marshall(testData, "/non-existing"+tmpFile)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

	outData := TestType{}
	err = marshaller.Unmarshall(tmpFile, &outData)
	if err != nil {
		t.Errorf("Unexpected error, got %s", err.Error())
	}

	if !reflect.DeepEqual(testData, outData) {
		t.Errorf("In and out does not match")
	}

	// Test Junk Content
	err = marshaller.Unmarshall(junkFile, &outData)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

	// Test Non Existing file
	err = marshaller.Unmarshall(junkFile+"-non-existing", &outData)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}
