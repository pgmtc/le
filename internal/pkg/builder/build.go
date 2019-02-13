package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/mholt/archiver"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pgmtc/orchard-cli/internal/pkg/docker"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var buildAction = common.RawAction{
	Handler: func(ctx common.Context, args ...string) error {
		noCache := false
		specDir := BUILDER_DIR

		for idx, arg := range args {
			if arg == "--nocache" {
				noCache = true
			}
			if arg == "--specdir" {
				if len(args) <= idx+1 {
					return fmt.Errorf("missing parameter for --specdir")
				}
				specDir = args[idx+1]
				ctx.Log.Debugf("Using %s as build spec dir\n", specDir)
			}
		}

		err, image, buildRoot, dockerFile, buildArgs := parseBuildProperties(specDir)
		if err != nil {
			return err
		}
		return buildImage(ctx, image, buildRoot, dockerFile, buildArgs, noCache)
	},
}

func parseBuildProperties(builderDir string) (resultErr error, image string, buildRoot string, dockerFile string, buildArgs []string) {
	// Try to read builder config
	configDirPath := common.ParsePath(builderDir)
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		resultErr = errors.Errorf("Unable to determine build configuration: %s", err.Error())
		return
	}

	bcnf := buildConfig{}
	bcnfPath := path.Join(builderDir, CONFIG_FILENAME)
	err := common.YamlUnmarshall(bcnfPath, &bcnf)
	if err != nil {
		resultErr = errors.Errorf("Unable to parse config file %s: %s", bcnfPath, err.Error())
		return
	}
	image = bcnf.Image
	buildRoot = common.ParsePath(bcnf.BuildRoot)
	dockerFile = common.ParsePath(bcnf.Dockerfile)
	buildArgs = bcnf.BuildArgs
	return
}

func mkContextTar(contextDir string, dockerFile string) (tarFile string, resultErr error) {
	tmpDir, resultErr := ioutil.TempDir("", "")
	if resultErr != nil {
		return
	}
	tarFile = tmpDir + "/docker-context.tar"
	resultErr = archiver.Archive([]string{contextDir + "/", dockerFile}, tarFile)
	return
}

func parseBuildArgs(buildArgs []string) (result map[string]*string) {
	result = map[string]*string{}
	for _, buildArg := range buildArgs {
		var argName, argValue string
		argSplit := strings.Split(buildArg, ":")
		switch len(argSplit) {
		case 0:
			argName = ""
			argValue = ""
		case 1:
			argName = argSplit[0]
			argValue = argSplit[0]
		default:
			argName = argSplit[0]
			argValue = argSplit[1]
		}
		argName = strings.Trim(argName, " ")
		argValue = strings.Trim(argValue, " ")
		if strings.HasPrefix(argValue, "$") {
			argValue = os.Getenv(argValue[1:])
		}
		if argName != "" {
			result[argName] = &argValue
		}
	}
	return
}

func buildImage(ctx common.Context, image string, buildRoot string, dockerFile string, buildArgs []string, noCache bool) error {
	log := ctx.Log
	if dockerFile == "" || image == "" || buildRoot == "" {
		return errors.Errorf("Missing parameters: image: %s, buildRoot: %s, dockerFile: %s", image, buildRoot, dockerFile)
	}
	log.Debugf("Building image %s'\n - Build Root: %s\n - Dockerfile: %s\n - No Cache: %s\n", image, buildRoot, dockerFile, noCache)

	log.Debugf("Creating context tar ... \n")
	contextTarFileName, returnError := mkContextTar(buildRoot, dockerFile)
	if returnError != nil {
		return returnError
	}
	defer os.Remove(contextTarFileName)
	log.Debugf("Context tar: %s\n", contextTarFileName)

	log.Debugf("Building docker context from %s\n", contextTarFileName)
	dockerBuildContext, returnError := os.Open(contextTarFileName)
	if returnError != nil {
		return returnError
	}
	defer dockerBuildContext.Close()

	cli := docker.DockerGetClient()
	args := parseBuildArgs(buildArgs)

	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     false,
		Tags:           []string{image},
		Dockerfile:     "Dockerfile",
		BuildArgs:      args,
		NoCache:        noCache,
	}

	log.Debugf("Starting docker build ...\n")
	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, options)
	if err != nil {
		log.Errorf("%s", err.Error())
	}
	log.Debugf("Finished with build\n")
	//defer buildResponse.Body.Close()

	//log.Debugf("********* %s **********\n", buildResponse.OSType)
	//_, err = io.Copy(os.Stdout, buildResponse.Body)
	//if err != nil {
	//	log.Fatal(err, " :unable to read image build response")
	//}

	d := json.NewDecoder(buildResponse.Body)

	type Event struct {
		Stream         string `json:"stream"`
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

		//log.Debugf("%+v\n", event)
		switch true {
		case event.Error != "":
			return errors.Errorf("\nbuild error: %s", event.Error)
		case event.Progress != "" || event.Status != "":
			log.Debugf("\r%s: %s", event.Status, event.Progress)
			if event.ProgressDetail.Current == 0 {
				log.Debugf("\n")
			}

		case strings.TrimSuffix(event.Stream, "\n") != "":
			log.Debugf("%s", event.Stream)
			if strings.HasPrefix(event.Stream, "Successfully built ") {
				// Fish for image id
				//imageId = strings.Replace(event.Stream, "Successfully built ", "", 1)
				//imageId = strings.TrimSuffix(imageId, "\n")
			}
		}
	}
	return nil
}
