package local

import (
	"github.com/golang/mock/gomock"
	"github.com/pgmtc/le/pkg/common"
	"github.com/pgmtc/le/pkg/local/mocks"
	"testing"
)

func setUp() (ctx common.Context, components []common.Component) {
	components = []common.Component{
		{
			Name:     "test-component",
			DockerId: "test-component",
			Image:    "bitnami/redis:latest",
		},
	}
	config := common.CreateMockConfig(components)
	log := common.ConsoleLogger{}
	ctx = common.Context{
		Log:    log,
		Config: config,
	}
	return
}

func Test_CreateAction(t *testing.T) {
	ctx, components := setUp()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRunner := mocks.NewMockRunner(mockCtrl)

	action := getComponentAction(mockRunner.Create)
	mockRunner.EXPECT().Create(ctx, components[0]).Times(1)
	err := action.Run(ctx, "test-component")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func Test_StartAction(t *testing.T) {
	ctx, components := setUp()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRunner := mocks.NewMockRunner(mockCtrl)

	action := getComponentAction(mockRunner.Start)
	mockRunner.EXPECT().Start(ctx, components[0]).Times(1)
	err := action.Run(ctx, "test-component")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func Test_StopAction(t *testing.T) {
	ctx, components := setUp()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRunner := mocks.NewMockRunner(mockCtrl)

	action := getComponentAction(mockRunner.Stop)
	mockRunner.EXPECT().Stop(ctx, components[0]).Times(1)
	err := action.Run(ctx, "test-component")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func Test_RemoveAction(t *testing.T) {
	ctx, components := setUp()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRunner := mocks.NewMockRunner(mockCtrl)

	action := getComponentAction(mockRunner.Remove)
	mockRunner.EXPECT().Remove(ctx, components[0]).Times(1)
	err := action.Run(ctx, "test-component")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func Test_PullAction(t *testing.T) {
	ctx, components := setUp()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRunner := mocks.NewMockRunner(mockCtrl)

	action := getComponentAction(mockRunner.Pull)
	mockRunner.EXPECT().Pull(ctx, components[0]).Times(1)
	err := action.Run(ctx, "test-component")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func Test_LogAction(t *testing.T) {
	ctx, components := setUp()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRunner := mocks.NewMockRunner(mockCtrl)

	noFollowAction := logsComponentAction(mockRunner, false)
	mockRunner.EXPECT().Logs(ctx, components[0], false).Times(1)
	err := noFollowAction.Run(ctx, "test-component")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	followLogAction := logsComponentAction(mockRunner, true)
	mockRunner.EXPECT().Logs(ctx, components[0], true).Times(1)
	err = followLogAction.Run(ctx, "test-component")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func Test_StatusAction(t *testing.T) {
	ctx, _ := setUp()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRunner := mocks.NewMockRunner(mockCtrl)

	// Plain status
	statusAction := getRawAction(mockRunner.Status)
	mockRunner.EXPECT().Status(ctx).Times(1)
	err := statusAction.Run(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	// Plain status
	statusFollowAction := getRawAction(mockRunner.Status)
	mockRunner.EXPECT().Status(ctx, "-v", "-f", "5").Times(1)
	err = statusFollowAction.Run(ctx, "-v", "-f", "5")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

//func Test_status(t *testing.T) {
//	ctx, runner := setUp()
//	statusAction := getRawAction(runner.Status)
//	pullAction := getComponentAction(runner.Pull)
//	createAction := getComponentAction(runner.Create)
//	startAction := getComponentAction(runner.Start)
//	stopAction := getComponentAction(runner.Stop)
//	removeAction := getComponentAction(runner.Remove)
//
//	COMPONENT_NAME := "test-status-component"
//	IMAGE_NAME := "nginx:alpine"
//	DOCKER_ID := "test-status-container"
//
//	var config = common.CreateMockConfig([]common.Component{
//		common.Component{
//			Name:          COMPONENT_NAME,
//			Image:         IMAGE_NAME,
//			DockerId:      DOCKER_ID,
//			ContainerPort: 80,
//			HostPort:      9998,
//			TestUrl:       "http://localhost:9998",
//		},
//	})
//
//	ctx = common.Context{
//		Log:    common.ConsoleLogger{},
//		Config: config,
//	}
//
//	removeAction.Run(ctx, COMPONENT_NAME) // Ignore error if it does not exist
//
//	// Test - no image present
//	if err := statusAction.Run(ctx); err != nil {
//		t.Errorf("Expected no error to be returned, but got %s", err.Error())
//	}
//
//	// Pull image and run status
//	if err := pullAction.Run(ctx, COMPONENT_NAME); err != nil {
//		t.Errorf("Expected no error to be returned, but got %s", err.Error())
//	}
//	defer removeAction.Run(ctx, COMPONENT_NAME)
//	if err := statusAction.Run(ctx); err != nil {
//		t.Errorf("Expected no error to be returned, but got %s", err.Error())
//	}
//
//	// Create container and run status
//	if err := createAction.Run(ctx, COMPONENT_NAME); err != nil {
//		t.Errorf("Expected no error to be returned, but got %s", err.Error())
//	}
//	defer removeAction.Run(ctx, COMPONENT_NAME)
//	if err := statusAction.Run(ctx); err != nil {
//		t.Errorf("Expected no error to be returned, but got %s", err.Error())
//	}
//
//	// Start container
//	if err := startAction.Run(ctx, COMPONENT_NAME); err != nil {
//		t.Errorf("Expected no error to be returned, but got %s", err.Error())
//	}
//	if err := statusAction.Run(ctx); err != nil {
//		t.Errorf("Expected no error to be returned, but got %s", err.Error())
//	}
//
//	// Stop container
//	if err := stopAction.Run(ctx, COMPONENT_NAME); err != nil {
//		t.Errorf("Expected no error to be returned, but got %s", err.Error())
//	}
//	if err := statusAction.Run(ctx); err != nil {
//		t.Errorf("Expected no error to be returned, but got %s", err.Error())
//	}
//
//	// Verbose and follow
//	if err := statusAction.Run(ctx, "-v", "-f", "5"); err != nil {
//		t.Errorf("Expected no error to be returned, but got %s", err.Error())
//	}
//
//	// Follow only
//	if err := statusAction.Run(ctx, "-f", "1"); err != nil {
//		t.Errorf("Expected no error to be returned, but got %s", err.Error())
//	}
//}
