package builder

import (
	"github.com/pgmtc/le/pkg/common"
	"github.com/pkg/errors"
	"os"
	"path"
	"strings"
)

const BUILDER_DIR = ".builder"
const CONFIG_FILENAME = "config.yaml"
const DEFAULT_DOCKERFILE = ".builder/Dockerfile"
const DEFAULT_BUILDROOT = ""

type buildConfig struct {
	Image      string
	BuildRoot  string
	Dockerfile string
	BuildArgs  []string
}

func initAction(fsHandler common.FsHandler, marshaller common.Marshaller) common.Action {
	return &common.RawAction{
		Handler: func(ctx common.Context, args ...string) error {
			log := ctx.Log
			log.Debugf("Init Action\n")

			configDirPath := common.ParsePath(BUILDER_DIR)
			if _, err := fsHandler.Stat(configDirPath); !os.IsNotExist(err) {
				return errors.Errorf("Directory %s already exists, please remove it first", configDirPath)
			}

			if err := fsHandler.MkdirAll(configDirPath, os.ModePerm); err != nil {
				return errors.Errorf("Error when creating config directory: %s", err.Error())
			}

			configPath := path.Join(configDirPath, CONFIG_FILENAME)

			bcnf := buildConfig{
				Image:      "my-image",
				BuildRoot:  DEFAULT_BUILDROOT,
				Dockerfile: DEFAULT_DOCKERFILE,
				BuildArgs:  []string{"build_arg_1:example_value"},
			}

			if err := marshaller.Marshall(bcnf, configPath); err != nil {
				return errors.Errorf("Error when writing build config: %s", err.Error())
			}

			// Create empty dockerfile
			dfPath := path.Join(configDirPath, strings.Replace(DEFAULT_DOCKERFILE, BUILDER_DIR, "", 1))
			if _, err := fsHandler.Create(dfPath); err != nil {
				return errors.Errorf("Error when writing empty Dockerfile: %s", err.Error())
			}

			return nil
		},
	}
}
