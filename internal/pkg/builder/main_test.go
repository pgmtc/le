package builder

import (
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"testing"
)

func TestParse(t *testing.T) {
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
			name: "successTest",
			args: args{
				[]string{"help"},
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

func Test_buildActionHandler(t *testing.T) {
	// This does not really test anything - just if it runs build - should fail as there is no component
	err := buildActionHandler(common.Component{}, common.HandlerArguments{})
	if err == nil {
		t.Errorf("Expected error, got nothing")
	}

}
