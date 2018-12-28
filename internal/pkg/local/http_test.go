package local

import (
	"fmt"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"net/http"
	"testing"
)

func startLocalServer(t *testing.T) {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/test-redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/test", 301)
	})
	http.HandleFunc("/test-error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error 500")
	})
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		t.Errorf("Error when starting local http server (for testing): %s", err.Error())
	}
}

func Test_isResponding(t *testing.T) {
	go startLocalServer(t) // Spins up local web server
	type args struct {
		cmp common.Component
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test-no-url",
			args: args{
				cmp: common.Component{
					Name: "testComponent",
				}},
			want: "",
		},
		{
			name: "test-success",
			args: args{
				cmp: common.Component{
					Name:    "testComponent",
					TestUrl: "http://localhost:9999/test",
				}},
			want: "200",
		},
		{
			name: "test-404",
			args: args{
				cmp: common.Component{
					Name:    "testComponent",
					TestUrl: "http://localhost:9999/non-existing",
				}},
			want: "404",
		},
		{
			name: "test-redirect",
			args: args{
				cmp: common.Component{
					Name:    "testComponent",
					TestUrl: "http://localhost:9999/test-redirect",
				}},
			want: "301",
		},
		{
			name: "test-fail",
			args: args{
				cmp: common.Component{
					Name:    "testComponent",
					TestUrl: "http://localhost:9999/test-error",
				}},
			want: "500",
		},
		{
			name: "test-fail",
			args: args{
				cmp: common.Component{
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
