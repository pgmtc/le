package docker

import (
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

// getAuthString returns authstring used by docker client based on component's dockerAuth property
func getAuthString(dockerAuth string) (authString string, resultErr error) {
	if dockerAuth != "" {
		commandParts := strings.Split(dockerAuth, " ")
		var cmd *exec.Cmd
		if len(commandParts) < 2 {
			cmd = exec.Command(commandParts[0])
		} else {
			cmd = exec.Command(commandParts[0], commandParts[1:]...)
		}

		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			resultErr = err
			return
		}
		authString, resultErr = parseLoginCmd(string(stdoutStderr))
	}
	return
}

// parseLoginCmd translates docker console client's login command to login string used by docker sdk
// example input docker login -u username -p password https://server-name
func parseLoginCmd(loginOutput string) (authString string, resultError error) {
	split := strings.Split(loginOutput, " ")
	if len(split) != 7 {
		resultError = errors.Errorf("Unexpected number of items in aws docker login command, got %d, expected %d", len(split), 7)
		return
	}

	authConfig := types.AuthConfig{
		Username:      split[3],
		Password:      split[5],
		ServerAddress: split[6],
	}
	encodedJSON, _ := json.Marshal(authConfig)
	authString = base64.URLEncoding.EncodeToString(encodedJSON)
	return
}
