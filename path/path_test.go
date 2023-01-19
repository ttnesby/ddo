//go:build unit

package path_test

import (
	"ddo/path"
	"testing"
)

const (
	bicepPath       = "./test/data/resourceGroup.bicep"
	nonExistingPath = "./n/a/nonExistingFile.bicep"
)

func TestHome(t *testing.T) {
	t.Parallel()

	_ = path.Home()
}

func TestHomeAbs(t *testing.T) {
	t.Parallel()

	got := path.HomeAbs("./")
	want := path.Home()

	if got != want {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestRepoRoot(t *testing.T) {
	t.Parallel()

	_ = path.RepoRoot()
}

func TestRepoAbs(t *testing.T) {
	t.Parallel()

	_ = path.RepoAbs(bicepPath)
}

func TestAbsExists(t *testing.T) {
	t.Parallel()

	if !path.AbsExists(path.RepoAbs(bicepPath)) {
		t.Error("want true for existing file, got false")
	}
}

func TestAbsExistsInvalid(t *testing.T) {
	t.Parallel()

	if path.AbsExists(path.RepoAbs(nonExistingPath)) {
		t.Error("want false for non-existing file, got true")
	}
}
