package docker

import (
	"encoding/base64"
	"os"
	"testing"
)

func Test_getAuthString(t *testing.T) {
	type args struct {
		dockerAuth string
	}

	tests := []struct {
		name           string
		args           args
		wantAuthString bool
		wantErr        bool
	}{
		{
			name: "eu-west-1",
			args: args{
				dockerAuth: "aws ecr get-login --no-include-email --region eu-west-1",
			},
			wantErr:        !(os.Getenv("SKIP_AWS_TESTING") == ""),
			wantAuthString: os.Getenv("SKIP_AWS_TESTING") == "",
		},
		{
			name: "unknown-repo",
			args: args{
				dockerAuth: "do-some-rubbish",
			},
			wantErr:        true,
			wantAuthString: false,
		},
		{
			name: "no-dockerAuth",
			args: args{
				dockerAuth: "",
			},
			wantErr:        false,
			wantAuthString: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAuthString, err := getAuthString(tt.args.dockerAuth)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAuthString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			decoded, err := base64.URLEncoding.DecodeString(gotAuthString)
			if err != nil {
				t.Errorf("Error when decoding: %s", err.Error())
			}
			decodedStr := string(decoded)
			if tt.wantAuthString == (decodedStr == "") {
				t.Errorf("getAuthStrig()\n - decodedStr = %v\n - wantAuthString = %v", decodedStr, tt.wantAuthString)
			}
		})
	}
}

func Test_parseLoginCmd(t *testing.T) {
	type args struct {
		loginOutput string
	}
	tests := []struct {
		name           string
		args           args
		wantAuthString string
		wantErr        bool
	}{
		{
			name: "testSuccess",
			args: args{
				loginOutput: "docker login -u username -p password https://server-name",
			},
			wantAuthString: "{\"username\":\"username\",\"password\":\"password\",\"serveraddress\":\"https://server-name\"}",
			wantErr:        false,
		},
		{
			name: "testFail",
			args: args{
				loginOutput: "some other unexpected return value",
			},
			wantAuthString: "",
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAuthString, err := parseLoginCmd(tt.args.loginOutput)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLoginCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			decodedAuthByte, _ := base64.URLEncoding.DecodeString(gotAuthString)
			decodedAuthString := string(decodedAuthByte)
			if decodedAuthString != tt.wantAuthString {
				t.Errorf("parseLoginCmd() = %v, want %v", gotAuthString, tt.wantAuthString)
			}
		})
	}
}
