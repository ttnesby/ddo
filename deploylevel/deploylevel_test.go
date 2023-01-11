package deploylevel_test

import (
	"ddo/deploylevel"
	"ddo/deployoperation"
	"ddo/reporoot"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

const (
	subId      = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	rgName     = "myrg"
	context    = "mycontext"
	deployment = "mycontext-895af37b7850639c7f2c"
	bicep      = "./test/data/test.bicep"
	json       = "./test/data/params.json"
)

func TestNewResourceGroup(t *testing.T) {
	t.Parallel()

	want := deploylevel.ResourceGroup{
		Deployment:        deployment,
		SubscriptionId:    uuid.MustParse(subId),
		ResourceGroupName: rgName,
		Template:          filepath.Join(reporoot.Get(), bicep),
	}

	got, err := deploylevel.NewResourceGroup(subId, rgName, bicep, context)
	if err != nil {
		t.Errorf("NewResourceGroup() returned error: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestNewResourceGroupInvalidUUIF(t *testing.T) {
	t.Parallel()

	if _, err := deploylevel.NewResourceGroup(
		"no uuid",
		rgName,
		bicep,
		context,
	); err == nil {
		t.Fatal("want error for invalid resource group, got nil")
	}
}

func TestNewResourceGroupInvalidTemplate(t *testing.T) {
	t.Parallel()

	if _, err := deploylevel.NewResourceGroup(
		subId,
		rgName,
		"no file",
		context,
	); err == nil {
		t.Fatal("want error for invalid resource group, got nil")
	}
}

func TestResourceGroupAZCmd(t *testing.T) {
	t.Parallel()

	rg, err := deploylevel.NewResourceGroup(subId, rgName, bicep, context)
	if err != nil {
		t.Errorf("NewResourceGroup() returned error: %v", err)
	}

	for _, op := range deployoperation.Operations() {
		got, err := rg.AZCmd(op, json)
		if err != nil {
			t.Errorf("AZCmd returned error: %v", err)
		}

		want := strings.Join(
			[]string{
				"az deployment group",
				string(op),
				"--name",
				deployment,
				"--subscription",
				subId,
				"--resource-group",
				rgName,
				"--template-file",
				filepath.Join(reporoot.Get(), bicep),
				"--out",
				"yaml",
				"--parameters",
				"@" + filepath.Join(reporoot.Get(), json),
			},
			" ")

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	}
}

func TestResourceGroupAZCmdInvalidParameterFile(t *testing.T) {
	t.Parallel()

	rg, err := deploylevel.NewResourceGroup(
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"myrg",
		bicep,
		"mycontext",
	)
	if err != nil {
		t.Errorf("NewResourceGroup() returned error: %v", err)
	}
	for _, op := range deployoperation.Operations() {
		if _, err := rg.AZCmd(op, "no param json"); err == nil {
			t.Fatal("want error for invalid AZCmd, got nil")
		}
	}
}
