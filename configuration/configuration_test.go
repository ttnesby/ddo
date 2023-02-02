package configuration_test

import (
	"ddo/alogger"
	cf "ddo/configuration"
	"ddo/path"
	"ddo/util"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/oklog/ulid/v2"
	"os"
	"path/filepath"
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

func TestConfigAsJson(t *testing.T) {
	t.Parallel()

	util.SkipIfCommandNotAvailable(t, requiredCommand)

	config := cf.New(rgConfigPath, []string{tagTenantNAVUtv})
	_, err := config.AsJson().Run()

	if err != nil {
		t.Errorf("could not extract config %v as json", err)
	}
}

func TestConfigAsYaml(t *testing.T) {
	t.Parallel()

	util.SkipIfCommandNotAvailable(t, requiredCommand)

	config := cf.New(rgConfigPath, []string{tagTenantNAVUtv})
	_, err := config.AsYaml().Run()

	if err != nil {
		t.Errorf("could not extract config %v as yaml", err)
	}
}

func TestConfigElementsAsJsonParametersTarget(t *testing.T) {
	t.Parallel()

	util.SkipIfCommandNotAvailable(t, requiredCommand)

	config := cf.New(rgConfigPath, []string{tagTenantNAVUtv})
	_, err := config.ElementsAsJson([]string{iElemParameters, iElemTarget}).Run()

	if err != nil {
		t.Errorf("could not extract config elements as json - %v", err)
	}
}

func TestConfigElementsAsTextRGName(t *testing.T) {
	t.Parallel()

	util.SkipIfCommandNotAvailable(t, requiredCommand)

	config := cf.New(rgConfigPath, []string{tagTenantNAVUtv})
	got, err := config.ElementsAsText([]string{iElemRGName}).Run()

	if err != nil {
		t.Errorf("could not extract config elements as text - %v", err)
	}

	if string(got) != iRgName {
		t.Errorf("want [%s], got [%s]", iRgName, got)
	}
}

func TestConfigElementsToTmpJsonFile(t *testing.T) {
	t.Parallel()

	util.SkipIfCommandNotAvailable(t, requiredCommand)

	absolutePath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("ddo.parameters.%s.json", ulid.Make().String()),
	)

	config := cf.New(rgConfigPath, []string{tagTenantNAVUtv})
	cmd := config.ElementsToTmpJsonFile(absolutePath, []string{iElemParameters})

	_, err := cmd.Run()

	if err != nil {
		t.Errorf("could not extract config elements to json file - %v", err)
	}

	if !path.AbsExists(absolutePath) {
		t.Errorf("Cannot find json file - %v", config)
	}
}
