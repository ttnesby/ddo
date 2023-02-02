package path

import (
	"ddo/alogger"
	"ddo/arg"
	"fmt"
	"github.com/oklog/ulid/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var l alogger.ALogger

func Init() {
	l = alogger.New(arg.InDebugMode())
}

// Home returns the absolute path of the home directory.
func Home() string {

	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("retrieving home path returned error: %v", err))
	}
	return strings.TrimRight(string(home), "\n")
}

// HomeAbs returns the absolute path of the given path relative to the home directory.
func HomeAbs(relativePath string) string {

	return filepath.Join(Home(), relativePath)
}

// RepoRoot returns the absolute path of the repo root.
func RepoRoot() string {

	cmd := []string{"git", "rev-parse", "--show-toplevel"}

	rr, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		panic(fmt.Errorf("retrieving repo root path [%v] returned error: %v", cmd, err))
	}
	return strings.TrimRight(string(rr), "\n")
}

// RepoAbs returns the absolute path of the given path relative to the repo root.
func RepoAbs(relativePath string) string {

	return filepath.Join(RepoRoot(), relativePath)
}

// AbsExists returns true if the given absolute path exists.
func AbsExists(absolutePath string) bool {
	_, err := os.Stat(absolutePath)
	return err == nil
}

func ContainerTmpJson() string {
	return filepath.Join("/tmp", fmt.Sprintf("ddo.parameters.%s.json", ulid.Make().String()))
}

func ActionSpecification() (path []string) {

	var foundPaths []string

	absToRelative := func(absPath, fileName string) (relativePath string) {
		//TODO fix err handling
		relativePath, _ = filepath.Rel(RepoRoot(), strings.TrimSuffix(absPath, fileName))
		return "./" + filepath.Clean(relativePath)
	}

	err := filepath.Walk(RepoRoot(), func(path string, info os.FileInfo, err error) error {

		if err != nil {
			l.Errorf("filepath.Walk() error: %v", err)
			return nil
		}

		if !info.IsDir() && info.Name() == "ddo.cue" {
			foundPaths = append(foundPaths, absToRelative(path, info.Name()))
		}

		return nil
	})

	if err != nil {
		l.Errorf("ActionSpecification() error: %v", err)
	}
	return foundPaths
}
