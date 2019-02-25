package common

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestOsFileSystemHandler_Stat(t *testing.T) {
	osHandler := OsFileSystemHandler{}
	_, err := osHandler.Stat("non-existing")
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

	tmpDir, _ := ioutil.TempDir("", "le-test-mock")
	defer os.RemoveAll(tmpDir)

	_, err = osHandler.Stat(tmpDir)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

}

func TestOsFileSystemHandler_Create(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "le-test-mock")
	defer os.RemoveAll(tmpDir)

	osHandler := OsFileSystemHandler{}
	_, err := osHandler.Create(filepath.Join(tmpDir, "some-file.txt"))
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	_, err = osHandler.Create("K:../../../../../../../../../../../../../../../../")
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}

func TestOsFileSystemHandler_MkdirAll(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "le-test-mock")
	defer os.RemoveAll(tmpDir)
	osHandler := OsFileSystemHandler{}
	err := osHandler.MkdirAll(filepath.Join(tmpDir, "some-dir"), os.ModePerm)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	err = osHandler.MkdirAll("../../../../../../../../../../../../../../../../K:", os.ModePerm)
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}
