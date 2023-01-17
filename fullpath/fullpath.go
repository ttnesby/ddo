package fullpath

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Get(path string) (string, error) {

	rr, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("retrieving repo root path returned error: %v", err)
	}

	fullPath := filepath.Join(strings.TrimRight(string(rr), "\n"), path)

	if _, err := os.Stat(fullPath); err != nil {
		return "", fmt.Errorf("checking path existence for %s returned error %v", fullPath, err)
	}

	return fullPath, nil
}
