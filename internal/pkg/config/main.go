package config

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/headzoo/surf"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Parse(args []string) error {
	actions := common.MakeActions()
	actions["status"] = status
	actions["init"] = initialize
	actions["create"] = create
	actions["switch"] = switchProfile
	actions["update-cli"] = updateCli
	return common.ParseParams(actions, args)
}

func status(args []string) error {
	color.HiWhite("Current profile: %s", common.CONFIG.Profile)
	color.HiWhite("Available profiles: %s", common.GetAvailableProfiles())
	if len(args) > 0 && args[0] == "-v" {
		// Verbose output
		s, _ := json.MarshalIndent(common.GetComponents(), "", "  ")
		color.White("Components: \n%s\n", s)
	} else {
		color.White("Components: (for more verbose output, add '-v' parameter)")
		for i, cmp := range common.GetComponents() {
			color.White("   %02d | Name: %s, DockerId: %s, Image: %s", i, cmp.Name, cmp.DockerId, cmp.Image)
		}
	}

	return nil
}

func switchProfile(args []string) error {
	if len(args) < 1 {
		return errors.Errorf("Missing parameter: profileName. Example: orchard config switch my-profile")
	}

	err := common.SwitchProfile(args[0])
	if err != nil {
		return errors.Errorf("Error when switching profile: %s", err.Error())
	}

	configFile, err := common.SaveConfig()
	if err != nil {
		return errors.Errorf("Erorr when saving config: %s", err.Error())
	}
	color.White("Successfully switched profile to %s. Changes written to %s", args[0], configFile)
	return nil
}

func initialize(args []string) error {
	common.SwitchCurrentProfile(common.DefaultLocalProfile())
	common.CONFIG.Profile = "default"

	fileName, err := common.SaveConfig()
	color.White("Config written to %s", fileName)

	fileName, err = common.SaveProfile("default", common.DefaultRemoteProfile())
	color.White("Profile written to %s", fileName)

	fileName, err = common.SaveProfile("local", common.DefaultLocalProfile())
	color.White("Profile written to %s", fileName)

	fileName, err = common.SaveProfile("remote", common.DefaultRemoteProfile())
	color.White("Profile written to %s", fileName)

	return err
}

func create(args []string) error {
	if len(args) < 1 {
		return errors.Errorf("Missing parameters: profileName [sourceProfile], examples:\n" +
			"    orchard config create my-new-profile\n" +
			"    orchard config create my-new-profile some-old-profile")
	}

	profileName := args[0]
	profile := common.DefaultLocalProfile()

	if len(args) > 1 {
		copyFromProfile, err := common.LoadProfile(args[1])
		if err != nil {
			return errors.Errorf("Error when loading profile %s: %s", args[1], err.Error())
		}
		profile = copyFromProfile
	}

	fileName, err := common.SaveProfile(profileName, profile)
	if err != nil {
		return errors.Errorf("Error when saving profile: %s", err.Error())
	}

	color.White("Successfully saved profile %s to %s", profileName, fileName)
	return nil
}

func updateCli(args []string) error {
	fmt.Printf("Determining latest release package .. ")
	url, fileName, err := getReleasePackagePath("x86", "macOS")
	if err != nil {
		return err
	}
	fmt.Printf("DONE\n")

	// Create temporary directory
	tmpDir, err := ioutil.TempDir("", "orchard-update")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	fmt.Printf("Downloading %s from %s ...", fileName, url)
	filePath := path.Join(tmpDir, fileName)
	out, err := os.Create(filePath)
	defer out.Close()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("DONE\n")
	fmt.Printf("Unzipping file %s ... ", fileName)
	unzippedFileName, err := unzipReleasePackage(filePath)
	if err != nil {
		return err
	}

	err = untar(unzippedFileName, tmpDir)
	if err != nil {
		return err
	}
	fmt.Printf("DONE\n")
	fmt.Printf("Installing orchard to %s .. ", common.CONFIG.BinLocation)
	newCliPath := path.Join(tmpDir, "orchard")
	err = installCli(newCliPath, common.CONFIG.BinLocation)
	if err != nil {
		return err
	}
	fmt.Printf("DONE\n")
	return nil

}

func unzipReleasePackage(fileName string) (resultFileName string, resultErr error) {
	gzipfile, err := os.Open(fileName)
	if err != nil {
		resultErr = err
		return
	}

	reader, err := gzip.NewReader(gzipfile)
	if err != nil {
		resultErr = err
		return
	}
	defer reader.Close()

	resultFileName = strings.TrimSuffix(fileName, ".gz")
	writer, err := os.Create(resultFileName)

	if err != nil {
		resultErr = err
		return
	}

	defer writer.Close()

	if _, err = io.Copy(writer, reader); err != nil {
		resultErr = err
		return
	}
	return
}

func getReleasePackagePath(arch string, os string) (urlString string, fileName string, resultErr error) {
	bow := surf.NewBrowser()
	err := bow.Open(common.CONFIG.ReleasesURL)
	if err != nil {
		resultErr = err
		return
	}

	for _, link := range bow.Links() {
		if strings.Contains(link.URL.Path, "releases/download") && strings.Contains(link.URL.Path, arch) && strings.Contains(link.URL.Path, os) {
			urlString = link.URL.String()
			pathSplit := strings.Split(link.URL.Path, "/")
			fileName = pathSplit[len(pathSplit)-1]
			return
		}
	}
	return
}

func untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func installCli(sourceFile string, targetFile string) error {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(targetFile, input, 755)
	if err != nil {
		return err
	}
	return nil
}
