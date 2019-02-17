package docker

import (
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

// getAuthString returns authstring used by docker client based on component's repository property
func getAuthString(repository string) (authString string, resultErr error) {
	var loginCmd string
	switch true {
	case strings.HasPrefix(repository, "ecr:"):
		loginCmd, resultErr = getEcrLoginCmd(repository)
	case repository == "":
		loginCmd = ""
	default:
		resultErr = errors.Errorf("Unknown repository type: %s", repository)
	}
	if resultErr != nil || loginCmd == "" {
		return
	}
	authString, resultErr = parseLoginCmd(loginCmd)
	return
}

// getEcrLoginCmd returns command in docker cli format
func getEcrLoginCmd(ecrLocation string) (loginCmd string, resultError error) {
	region := strings.TrimPrefix(ecrLocation, "ecr:")
	var cmdResults []byte
	cmdResults, resultError = exec.Command("aws", "ecr", "get-login", "--no-include-email", "--region", region).Output()
	if resultError != nil {
		return
	}
	loginCmd = string(cmdResults)
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
