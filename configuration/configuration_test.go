//go:build unit

package configuration_test

import (
	"ddo/alogger"
	cf "ddo/configuration"
	"ddo/path"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

const (
	tagTenantNAVUtv   = "tenant=navutv"
	tagSomething      = "something=2"
	rgConfigPath      = "./test/infrastructure/resourceGroup"
	elemParameters    = "parameters"
	elemTarget        = "target"
	elemRGName        = "#name"
	rgName            = "container-registry"
	invalidConfigPath = "./n/a"
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

	config := cf.New(invalidConfigPath, nil)
	_, err := config.AsJson()

	if err == nil {
		t.Error("want error for invalid config, got nil")
	}
}

func TestConfigAsJson(t *testing.T) {
	t.Parallel()

	config := cf.New(rgConfigPath, []string{tagTenantNAVUtv})
	_, err := config.AsJson()

	if err != nil {
		t.Errorf("could not extract config %v as json", err)
	}
}

func TestConfigAsYaml(t *testing.T) {
	t.Parallel()

	config := cf.New(rgConfigPath, []string{tagTenantNAVUtv})
	_, err := config.AsYaml()

	if err != nil {
		t.Errorf("could not extract config %v as yaml", err)
	}
}

func TestConfigElementsAsJsonParametersTarget(t *testing.T) {
	t.Parallel()

	config := cf.New(rgConfigPath, []string{tagTenantNAVUtv})
	_, err := config.ElementsAsJson([]string{elemParameters, elemTarget})

	if err != nil {
		t.Errorf("could not extract config elements as json - %v", err)
	}
}

func TestConfigElementsAsTextRGName(t *testing.T) {
	t.Parallel()

	config := cf.New(rgConfigPath, []string{tagTenantNAVUtv})
	got, err := config.ElementsAsText([]string{elemRGName})

	if err != nil {
		t.Errorf("could not extract config elements as text - %v", err)
	}

	if string(got) != rgName {
		t.Errorf("want [%s], got [%s]", rgName, got)
	}
}

func TestConfigElementsToTmpJsonFile(t *testing.T) {
	t.Parallel()

	config := cf.New(rgConfigPath, []string{tagTenantNAVUtv})
	got, err := config.ElementsToTmpJsonFile([]string{elemParameters})

	if err != nil {
		t.Errorf("could not extract config elements as text - %v", err)
	}

	if !path.AbsExists(got) {
		t.Errorf("Cannot find exported file - %v", config)
	}
}
