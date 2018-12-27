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
			args: args{
				cmp: Component{
					Name: "testComponent",
				}},
			want: "",
		},
		{
			name: "test-success",
			args: args{
				cmp: Component{
					Name:    "testComponent",
					TestUrl: "http://www.praguematica.co.uk/",
				}},
			want: "200",
		},
		{
			name: "test-redirect",
			args: args{
				cmp: Component{
					Name:    "testComponent",
					TestUrl: "http://praguematica.co.uk/",
				}},
			want: "301",
		},
		{
			name: "test-fail",
			args: args{
				cmp: Component{
					Name:    "testComponent",
					TestUrl: "http://non-existing-url/",
				}},
			want: "ERR",
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
