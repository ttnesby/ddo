package cuecli_test

import (
	"ddo/alogger"
	cf "ddo/cuecli"
	"ddo/util"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

const (
	tagTenantNAVUtv    = "tenant=navutv"
	tagSomething       = "something=2"
	rgConfigPath       = "./test/infrastructure/resourceGroup"
	iElemParameters    = "parameters"
	iElemTarget        = "target"
	iElemRGName        = "#name"
	iRgName            = "container-registry"
	iInvalidConfigPath = "./n/a"
)

const (
	requiredCommand = "cue"
)

func TestMain(m *testing.M) {
	//setup
	alogger.Disable()
	code := m.Run()
	//shutdown
	os.Exit(code)
}

func TestConfigWithTags(t *testing.T) {
	t.Parallel()

	want := cf.CueCli{"cue", "export", rgConfigPath, "-t", tagTenantNAVUtv, "-t", tagSomething}
	got := cf.New(rgConfigPath, []string{tagTenantNAVUtv, tagSomething})

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestConfigWithoutTags(t *testing.T) {
	t.Parallel()

	want := cf.CueCli{"cue", "export", rgConfigPath}
	got := cf.New(rgConfigPath, nil)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestConfigInvalid(t *testing.T) {
	t.Parallel()

	util.SkipIfCommandNotAvailable(t, requiredCommand)

	config := cf.New(iInvalidConfigPath, nil)
	_, err := config.AsJson().Run()

	if err == nil {
		t.Error("want error for invalid config, got nil")
	}
}
