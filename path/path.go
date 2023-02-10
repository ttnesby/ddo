package path

import (
	"ddo/alogger"
	"ddo/arg"
	"ddo/util"
	"fmt"
	"github.com/oklog/ulid/v2"
	"os"
	"path/filepath"
	"strings"
)

const (
	DotAzure          = ".azure"
	ContainerRef      = "docker.io/ttnesby/azbicue:latest"
	ContainerDotAzure = "/root/" + DotAzure
	ContainerUser     = "/root"
	ContainerRepoRoot = "/rr"
)

var l alogger.ALogger

func Init() {
	l = alogger.New(arg.InDebugMode())
}

func home() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("retrieving home path returned error: %v", err))
	}
	return home
}

func HostDotAzure() string {
	if util.InDockerContainer() {
		return filepath.Join(ContainerUser, DotAzure)
	}
	return filepath.Join(home(), DotAzure)
}

func HostRepoRoot() string {
	if util.InDockerContainer() {
		return ContainerRepoRoot
	}
	return "."
}

func ContainerTmpJson() (tmpJson string) {
	tmpJson = filepath.Join("/tmp", fmt.Sprintf("ddo.parameters.%s.json", ulid.Make().String()))
	l.Debugf("Creating temporary json-file %s", tmpJson)
	return tmpJson
}

func clean(p, f string) string {
	relPath := func() string {
		if filepath.IsAbs(p) {
			rp, _ := filepath.Rel(ContainerRepoRoot, strings.TrimSuffix(p, string(os.PathSeparator)+f))
			return rp
		}
		return strings.TrimSuffix(p, string(os.PathSeparator)+f)
	}()
	return "." + string(os.PathSeparator) + relPath
}

func ActionSpecification() (path []string) {

	var foundPaths []string

	err := filepath.Walk(HostRepoRoot(), func(path string, info os.FileInfo, err error) error {

		if err != nil {
			l.Errorf("filepath.Walk() error: %v", err)
			return nil
		}

		if !info.IsDir() && info.Name() == "ddo.cue" {
			foundPaths = append(foundPaths, clean(path, info.Name()))
		}

		return nil
	})

	if err != nil {
		l.Errorf("ActionSpecification() error: %v", err)
	}
	return foundPaths
}
