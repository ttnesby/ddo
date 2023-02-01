package path_test

import (
	"ddo/path"
	"ddo/util"
	"testing"
)

const (
	bicepPath       = "./test/data/resourceGroup.bicep"
	nonExistingPath = "./n/a/nonExistingFile.bicep"
)

const (
	requiredCommand = "git"
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

	util.SkipIfCommandNotAvailable(t, requiredCommand)

	_ = path.RepoRoot()
}

func TestRepoAbs(t *testing.T) {
	t.Parallel()

	util.SkipIfCommandNotAvailable(t, requiredCommand)

	_ = path.RepoAbs(bicepPath)
}

func TestAbsExistsRepo(t *testing.T) {
	t.Parallel()

	util.SkipIfCommandNotAvailable(t, requiredCommand)

	if !path.AbsExists(path.RepoAbs(bicepPath)) {
		t.Error("want true for existing file, got false")
	}
}

func TestAbsExistsInvalidRepo(t *testing.T) {
	t.Parallel()

	util.SkipIfCommandNotAvailable(t, requiredCommand)

	if path.AbsExists(path.RepoAbs(nonExistingPath)) {
		t.Error("want false for non-existing file, got true")
	}
}
