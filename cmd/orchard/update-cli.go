package main

import (
	"archive/tar"
	"compress/gzip"
	"github.com/headzoo/surf"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var updateCliAction = common.RawAction{
	Handler: func(ctx common.Context, args ...string) error {
		log := ctx.Log
		config := ctx.Config
		log.Debugf("Determining latest release package .. ")
		url, fileName, err := getReleasePackagePath(config.Config().ReleasesURL, "x86", "macOS")
		if err != nil {
			return err
		}
		log.Debugf("DONE\n")

		// Create temporary directory
		tmpDir, _ := ioutil.TempDir("", "orchard-update")
		defer os.RemoveAll(tmpDir) // clean up

		log.Debugf("Downloading %s from %s ...", fileName, url)
		filePath := path.Join(tmpDir, fileName)
		out, err := os.Create(filePath)
		defer out.Close()
		resp, err := http.Get(url)
		defer resp.Body.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
		log.Debugf("DONE\n")
		log.Debugf("Unzipping file %s ... ", fileName)
		unzippedFileName, err := unzipReleasePackage(filePath)
		if err != nil {
			return err
		}

		err = untar(unzippedFileName, tmpDir)
		if err != nil {
			return err
		}
		log.Debugf("DONE\n")
		log.Debugf("Installing orchard to %s .. ", config.Config().BinLocation)
		newCliPath := path.Join(tmpDir, "orchard")
		err = installCli(newCliPath, config.Config().BinLocation)
		if err != nil {
			return err
		}
		log.Debugf("DONE\n")
		return nil
	},
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

func getReleasePackagePath(releasesUrl string, arch string, os string) (urlString string, fileName string, resultErr error) {
	bow := surf.NewBrowser()
	err := bow.Open(releasesUrl)
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
