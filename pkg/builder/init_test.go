package builder

import (
	"github.com/golang/mock/gomock"
	"github.com/pgmtc/le/pkg/common"
	"github.com/pgmtc/le/pkg/common/mocks"
	"github.com/pkg/errors"
	"os"
	"testing"
)

func cleanup() {
	if _, err := os.Stat(BUILDER_DIR); !os.IsNotExist(err) {
		os.RemoveAll(BUILDER_DIR)
	}
}

func setUpInit(t *testing.T) (mockCtrl *gomock.Controller, mockFsHandler *mocks.MockFsHandler, mockMarshaller *mocks.MockMarshaller, action common.Action) {
	mockCtrl = gomock.NewController(t)
	mockFsHandler = mocks.NewMockFsHandler(mockCtrl)
	mockMarshaller = mocks.NewMockMarshaller(mockCtrl)
	action = initAction(mockFsHandler, mockMarshaller)
	return
}

func Test_initAction_success(t *testing.T) {
	mockCtrl, mockFsHandler, mockMarshaller, action := setUpInit(t)
	defer mockCtrl.Finish()

	// Success scenario
	mockFsHandler.EXPECT().Stat(gomock.Any()).Return(nil, os.ErrNotExist)
	mockFsHandler.EXPECT().MkdirAll(gomock.Any(), gomock.Any()).Return(nil)
	mockMarshaller.EXPECT().Marshall(gomock.Any(), gomock.Any()).Return(nil)
	mockFsHandler.EXPECT().Create(gomock.Any()).Return(nil, nil)
	err := action.Run(common.Context{
		Config: common.CreateMockConfig([]common.Component{}),
		Log:    common.ConsoleLogger{},
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func Test_initAction_dirExists(t *testing.T) {
	mockCtrl, mockFsHandler, _, action := setUpInit(t)
	defer mockCtrl.Finish()

	// Success scenario
	mockFsHandler.EXPECT().Stat(gomock.Any()).Return(nil, nil)
	err := action.Run(common.Context{
		Config: common.CreateMockConfig([]common.Component{}),
		Log:    common.ConsoleLogger{},
	})
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}

func Test_initAction_createFailure(t *testing.T) {
	mockCtrl, mockFsHandler, _, action := setUpInit(t)
	defer mockCtrl.Finish()

	// Success scenario
	mockFsHandler.EXPECT().Stat(gomock.Any()).Return(nil, os.ErrNotExist)
	mockFsHandler.EXPECT().MkdirAll(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
	err := action.Run(common.Context{
		Config: common.CreateMockConfig([]common.Component{}),
		Log:    common.ConsoleLogger{},
	})
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}

func Test_initAction_marshallFailure(t *testing.T) {
	mockCtrl, mockFsHandler, mockMarshaller, action := setUpInit(t)
	defer mockCtrl.Finish()

	// Success scenario
	mockFsHandler.EXPECT().Stat(gomock.Any()).Return(nil, os.ErrNotExist)
	mockFsHandler.EXPECT().MkdirAll(gomock.Any(), gomock.Any()).Return(nil)
	mockMarshaller.EXPECT().Marshall(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
	err := action.Run(common.Context{
		Config: common.CreateMockConfig([]common.Component{}),
		Log:    common.ConsoleLogger{},
	})
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}

func Test_initAction_writeFailure(t *testing.T) {
	mockCtrl, mockFsHandler, mockMarshaller, action := setUpInit(t)
	defer mockCtrl.Finish()

	// Success scenario
	mockFsHandler.EXPECT().Stat(gomock.Any()).Return(nil, os.ErrNotExist)
	mockFsHandler.EXPECT().MkdirAll(gomock.Any(), gomock.Any()).Return(nil)
	mockMarshaller.EXPECT().Marshall(gomock.Any(), gomock.Any()).Return(nil)
	mockFsHandler.EXPECT().Create(gomock.Any()).Return(nil, errors.New("mock error"))
	err := action.Run(common.Context{
		Config: common.CreateMockConfig([]common.Component{}),
		Log:    common.ConsoleLogger{},
	})
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}
}
