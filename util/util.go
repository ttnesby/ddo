package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)

const (
	requiredCommand = "az"
	navutv          = "82bdf6c1-3e56-4a5e-8c50-c331165e0192"
)

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func SkipIfCommandNotAvailable(t *testing.T, cmd string) {
	if !CommandExists(cmd) {
		t.Skipf("%s not installed", cmd)
	}
}

func SkipIfNotLoggedIntoNAVUTV(t *testing.T) {
	if !CommandExists(requiredCommand) {
		t.Skip(fmt.Sprintf("%s not installed", requiredCommand))
	}

	tenantIdCmd := []string{"az", "account", "show", "--query", "homeTenantId", "--output", "tsv"}
	tenantId, err := exec.Command(tenantIdCmd[0], tenantIdCmd[1:]...).CombinedOutput()
	if err != nil {
		t.Skip(fmt.Sprintf("most likely, not logged in, executing %v: %v", tenantIdCmd, err))
	}

	if string(bytes.TrimRight(tenantId, "\r\n")) != navutv {
		t.Skip("not logged into tenant NAVUTV")
	}
}
