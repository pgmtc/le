package source

import "testing"

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
			name: "test-fail",
			args: args {
				[]string{"nonExistingAction", "param2", "param3"},
				},
			wantErr: true,
		},
		{
			name: "test-pull",
			args: args {
				[]string{"pull", "param2", "param3"},
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
