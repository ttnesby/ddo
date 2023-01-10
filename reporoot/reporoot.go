package reporoot

import (
	"fmt"
	"os/exec"
	"strings"
)

func Get() string {

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	rr, err := cmd.Output()

	if err != nil {
		panic(fmt.Errorf("cmd.Output() returned error: %v", err))
	}

	return strings.TrimRight(string(rr), "\n")
}
