package deploylevel_test

import (
	"ddo/category"
	"ddo/deploylevel"
	"ddo/reporoot"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestNewResourceGroup(t *testing.T) {
	t.Parallel()

	bicep := "/test/data/test.bicep"
	bicepFile := filepath.Join(reporoot.Get(), bicep)

	if _, err := os.Stat(bicepFile); err != nil {
		t.Errorf("os.Stat() returned error: %v", err)
	}

	want := deploylevel.ResourceGroup{
		Deployment:        "mycontext-895af37b7850639c7f2c",
		SubscriptionId:    uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
		ResourceGroupName: "myrg",
		Template:          bicepFile,
	}

	got, err := deploylevel.NewResourceGroup("6ba7b810-9dad-11d1-80b4-00c04fd430c8", "myrg", bicep, "mycontext")
	if err != nil {
		t.Errorf("NewResourceGroup() returned error: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(want, got))
	}
}
func TestNewResourceGroupInvalid(t *testing.T) {
	t.Parallel()

	_, err := deploylevel.NewResourceGroup("no uuid", "myrg", "nofile", "mycontext")
	if err == nil {
		t.Fatal("want error for invalid resource group, got nil")
	}

	_, err = deploylevel.NewResourceGroup("6ba7b810-9dad-11d1-80b4-00c04fd430c8", "myrg", "nofile", "mycontext")
	if err == nil {
		t.Fatal("want error for invalid resource group, got nil")
	}
}

func TestResourceGroupAZCmd(t *testing.T) {
	t.Parallel()

	bicep := "./test/data/test.bicep"
	bicepFile := filepath.Join(reporoot.Get(), bicep)

	json := "./test/data/params.json"
	jsonFile := filepath.Join(reporoot.Get(), json)

	rg, err := deploylevel.NewResourceGroup("6ba7b810-9dad-11d1-80b4-00c04fd430c8", "myrg", bicep, "mycontext")
	if err != nil {
		t.Errorf("NewResourceGroup() returned error: %v", err)
	}

	got, err := rg.AZCmd(category.Validate, json)
	if err != nil {
		t.Errorf("AZCmd returned error: %v", err)
	}

	want := strings.Join(
		[]string{
			"az deployment group",
			"validate",
			"--name",
			"mycontext-895af37b7850639c7f2c",
			"--subscription",
			"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			"--resource-group",
			"myrg",
			"--template-file",
			bicepFile,
			"--out",
			"yaml",
			"--parameters",
			"@" + jsonFile,
		},
		" ")

	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}

}
