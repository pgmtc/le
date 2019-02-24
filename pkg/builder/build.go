package builder

import (
	"fmt"
	"github.com/pgmtc/le/pkg/common"
	"github.com/pkg/errors"
	"os"
	"path"
)

func getBuildAction(builder Builder) common.Action {
	return &common.RawAction{
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
			return builder.BuildImage(ctx, image, buildRoot, dockerFile, buildArgs, noCache)
		},
	}
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
