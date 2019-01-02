package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/fatih/color"
	"github.com/jhoonb/archivex"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pgmtc/orchard-cli/internal/pkg/local"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func relativeOrAbsolute(path string) string {
	if !strings.HasPrefix(path, "/") {
		currentDir, _ := os.Getwd()
		return currentDir + "/" + path
	}
	return path
}

func findMsvcJar(path string) (fileName string, returnError error) {
	files, returnError := filepath.Glob(path + "/target/*.jar")
	if returnError != nil {
		return
	}

	if len(files) > 1 {
		returnError = errors.Errorf("Unexpected number of jar files found - expecting only one in '%s'", path + "/target/")
		return
	}

	for _, file := range files {
		fileName = strings.Replace(file, path + "/", "", 1)
	}
	return
}

func build(component common.Component, handlerArguments common.HandlerArguments) (imageId string, returnError error) {
	fmt.Printf("Build an image, %s\n", component.Name)
	buildRoot := relativeOrAbsolute(component.BuildRoot)
	dockerFile := relativeOrAbsolute(component.DockerFile)
	contextTarFileName, returnError := mkContextTar(buildRoot, dockerFile)
	if returnError != nil {
		return
	}
	defer os.Remove(contextTarFileName)

	dockerBuildContext, returnError := os.Open(contextTarFileName)
	if returnError != nil {
		return
	}
	defer dockerBuildContext.Close()
	cli := local.DockerGetClient()

	artifactoryPassword := os.Getenv("ARTIFACTORY_PASSWORD")

	jarFile, err := findMsvcJar(buildRoot)
	if err != nil {
		returnError = errors.Errorf("problem determining jar file for msvc: %s", err.Error())
		return
	}

	if jarFile != "" {
		color.Blue("JAR_FILE used: %s", jarFile)
	}

	args := map[string]*string{
		"mvn_password":     &artifactoryPassword,
		"JAR_FILE": 		&jarFile,
	}


	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     false,
		Tags:			[]string {component.Image},
		Dockerfile:     "Dockerfile",
		BuildArgs: args,
	}
	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, options)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	defer buildResponse.Body.Close()

	//fmt.Printf("********* %s **********\n", buildResponse.OSType)
	//_, err = io.Copy(os.Stdout, buildResponse.Body)
	//if err != nil {
	//	log.Fatal(err, " :unable to read image build response")
	//}

	d := json.NewDecoder(buildResponse.Body)

	type Event struct {
		Stream 		   string `json:"stream"`
		Status         string `json:"status"`
		Error          string `json:"error"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"progressDetail"`
	}

	var event *Event
	for {
		if err := d.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		//fmt.Printf("%+v\n", event)
		switch true {
		case event.Error != "":
			returnError = errors.Errorf("build error: %s", event.Error)
			return
		case event.Progress != "" || event.Status != "":
			fmt.Printf(color.MagentaString("\r%s: %s", event.Status, event.Progress))
			if event.ProgressDetail.Current == 0 {
				fmt.Println()
			}

		case strings.TrimSuffix(event.Stream, "\n") != "":
			fmt.Printf(color.MagentaString("%s", event.Stream))
			if strings.HasPrefix(event.Stream, "Successfully built ") {
				// Fish for image id
				imageId = strings.Replace(event.Stream, "Successfully built ", "", 1)
				imageId = strings.TrimSuffix(imageId, "\n")
			}
		}}

	return
}

func mkContextTar(contextDir string, dockerFile string) (string, error) {
	if _, err := os.Stat(contextDir); os.IsNotExist(err) {
		return "", errors.Errorf("Context directory '%s' does not exist", contextDir)
	}
	// Create temporary file
	tmpfile, err := ioutil.TempFile("", "context-")
	if err != nil {
		return "", errors.Errorf("problem when creating temporary file: %s", err.Error())
	}

	// Create tar
	tarFileName := tmpfile.Name() + ".tar"
	tar := new(archivex.TarFile)
	if err := tar.Create(tarFileName); err != nil {
		return "", errors.Errorf("error creating tar file: %s", err.Error())
	}

	fr, err := os.Open(dockerFile)
	if err != nil {
		return "", errors.Errorf("error when reading dockerfile: %s", err.Error())
	}

	tar.AddAll(contextDir, false)
	tar.Add("Dockerfile", fr, nil)
	tar.Close()

	return tar.Name, nil
}

func RemoveImage(imageId string) {
	local.DockerGetClient().ImageRemove(context.Background(), imageId, types.ImageRemoveOptions{
		Force:true,
		PruneChildren:true,
	})
}
