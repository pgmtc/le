package repo

import (
	"github.com/pgmtc/le/pkg/common"
	"testing"
)

func TestModule_GetActions(t *testing.T) {
	mods := Module{}.GetActions()
	if len(mods) == 0 {
		t.Errorf("Unexpected number of modules: %d", len(mods))
	}
}

func Test_urlAction(t *testing.T) {
	logger := &common.StringLogger{}
	cnf := common.CreateMockConfig([]common.Component{{Name: "test-component"}})

	action := &urlAction
	err := action.Run(common.Context{
		Config: cnf,
		Log:    logger,
	}, "test-repository")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if len(logger.InfoMessages) != 1 {
		t.Errorf("Unexpected length of the log. Expected 1, got %d", len(logger.InfoMessages))
	}

	expected := "git clone " + cnf.Config().RepositoryPrefix + "test-repository\n"
	if logger.InfoMessages[0] != expected {
		t.Errorf("Unexpected result. \nexp: %sgot: %s", expected, logger.InfoMessages[0])
	}
}
