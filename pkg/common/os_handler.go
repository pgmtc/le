package common

import "os"

//go:generate mockgen -destination=./mocks/mock_fshandler.go -package=mocks github.com/pgmtc/le/pkg/common FsHandler

type FsHandler interface {
	MkdirAll(path string, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	Create(name string) (*os.File, error)
}

type OsFileSystemHandler struct{}

func (OsFileSystemHandler) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OsFileSystemHandler) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (OsFileSystemHandler) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
