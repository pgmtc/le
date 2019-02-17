package docker

import (
	"encoding/base64"
	"os"
	"regexp"
	"strings"
	"testing"
)

func Test_getAuthString(t *testing.T) {
	type args struct {
		repository string
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
				repository: "ecr:eu-west-1",
			},
			wantErr:        !(os.Getenv("SKIP_AWS_TESTING") == ""),
			wantAuthString: true,
		},
		{
			name: "unknown-repo",
			args: args{
				repository: "some-unknown-repo",
			},
			wantErr:        true,
			wantAuthString: false,
		},
		{
			name: "no-repository",
			args: args{
				repository: "",
			},
			wantErr:        false,
			wantAuthString: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAuthString, err := getAuthString(tt.args.repository)
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

func Test_getEcrLoginCmd(t *testing.T) {
	if os.Getenv("SKIP_AWS_TESTING") != "" {
		t.Skipf("SKIP_AWS_TESTING present, skipping")
	}
	type args struct {
		ecrLocation string
	}
	tests := []struct {
		name        string
		args        args
		wantToMatch string
		wantErr     bool
	}{
		{
			name: "test us-east-1",
			args: args{
				ecrLocation: "ecr:us-east-1",
			},
			wantErr: false,
		},
		{
			name: "test eu-west-1",
			args: args{
				ecrLocation: "ecr:eu-west-1",
			},
			wantErr: false,
		},
		{
			name: "test moon (non existing region)",
			args: args{
				ecrLocation: "ecr:moon",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLoginCmd, err := getEcrLoginCmd(tt.args.ecrLocation)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEcrLoginCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				region := strings.TrimPrefix(tt.args.ecrLocation, "ecr:")
				regex := "docker login -u AWS -p [^\\s]+ https://[^\\s\\\\.]+.dkr.ecr." + region + ".amazonaws.com"
				r, _ := regexp.Compile(regex)
				if !r.MatchString(gotLoginCmd) {
					t.Errorf("Unexpected result\n - %s - not matching regexp: %s", gotLoginCmd, regex)
				}
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
