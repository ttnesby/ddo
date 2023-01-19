package path

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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
