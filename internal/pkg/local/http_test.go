package local

import "testing"

func Test_isResponding(t *testing.T) {
	type args struct {
		cmp Component
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test-no-url",
			args: args {
				cmp: Component{
					name: "testComponent",
				}},
			want: "",
		},
		{
			name: "test-200",
			args: args {
				cmp: Component{
					name: "testComponent",
					testUrl: "http://www.praguematica.co.uk/",
				}},
			want: "200",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isResponding(tt.args.cmp); got != tt.want {
				t.Errorf("isResponding() = %v, want %v", got, tt.want)
			}
		})
	}
}
