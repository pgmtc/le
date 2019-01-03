package config

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"io/ioutil"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	var origConfigLocation, tmpDir string
	setUp := func() {
		origConfigLocation = common.CONFIG_LOCATION
		tmpDir, _ = ioutil.TempDir("", "orchard-test-config-mock")
		common.CONFIG_LOCATION = tmpDir
	}

	rollBack := func() {
		common.CONFIG_LOCATION = origConfigLocation
		os.RemoveAll(tmpDir)
	}

	setUp()
	defer rollBack()

	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "failTest",
			args: args{
				[]string{"nonExistingAction"},
			},
			wantErr: true,
		},
		{
			name: "initTest",
			args: args{
				[]string{"init"},
			},
			wantErr: false,
		},
		{
			name: "createTest",
			args: args{
				[]string{"create", "profile1"},
			},
			wantErr: false,
		},
		{
			name: "createFromTest",
			args: args{
				[]string{"create", "profile2", "profile1"},
			},
			wantErr: false,
		},
		{
			name: "createFromFailTest",
			args: args{
				[]string{"create", "profile3", "profile0"},
			},
			wantErr: true,
		},
		{
			name: "createFailNoParams",
			args: args{
				[]string{"create"},
			},
			wantErr: true,
		},
		{
			name: "switchTest",
			args: args{
				[]string{"switch", "profile1"},
			},
			wantErr: false,
		},
		{
			name: "switchTestFailNonExisting",
			args: args{
				[]string{"switch", "profile0"},
			},
			wantErr: true,
		},
		{
			name: "switchTestFailNoParams",
			args: args{
				[]string{"switch"},
			},
			wantErr: true,
		},
		{
			name: "statusTest",
			args: args{
				[]string{"status"},
			},
			wantErr: false,
		},
		{
			name: "statusVerboseTest",
			args: args{
				[]string{"status", "-v"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Parse(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initialize(t *testing.T) {
	err := initialize([]string{})
	if err != nil {
		t.Errorf("Unexpected error returned: %s", err.Error())
	}
}
