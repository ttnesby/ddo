//go:build unit

package fullpath_test

import (
	"ddo/fullpath"
	"testing"
)

const (
	bicepPath       = "./test/data/resourceGroup.bicep"
	nonExistingPath = "./n/a/nonExistingFile.bicep"
)

func TestGet(t *testing.T) {
	t.Parallel()

	if _, err := fullpath.Get(bicepPath); err != nil {
		t.Error(err)
	}
}

func TestGetNonExisting(t *testing.T) {
	t.Parallel()

	if _, err := fullpath.Get(nonExistingPath); err == nil {
		t.Error("want error for non-existing file, got nil")
	}
}
