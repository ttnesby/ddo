package deploy_test

import (
	"ddo/deploy"
	dl "ddo/deploy/level"
	do "ddo/deploy/operation"
	rr "ddo/reporoot"
	"path/filepath"
	"strings"
	"testing"
)

const (
	location   = "norwayEast"
	mgId       = "6ba7b810-9dad-11d1-80b4-00c04fd43022"
	subId      = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	rgName     = "myRg"
	context    = "myContext"
	deployment = "myContext-35a2b15e5a8d96e4d362"
	bicep      = "./test/data/test.bicep"
	json       = "./test/data/params.json"
)

func TestNewResourceGroup(t *testing.T) {
	t.Parallel()

	for _, op := range do.Operations() {
		got, err := deploy.New(dl.ResourceGroup, op, context, subId, rgName, bicep, json)
		if err != nil {
			t.Errorf("New() returned error: %v", err)
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
				filepath.Join(rr.Get(), bicep),
				"--out",
				"yaml",
				"--parameters",
				"@" + filepath.Join(rr.Get(), json),
			},
			" ")

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	}
}

func TestNewSubscription(t *testing.T) {
	t.Parallel()

	for _, op := range do.Operations() {
		got, err := deploy.New(dl.Subscription, op, context, subId, location, bicep, json)
		if err != nil {
			t.Errorf("New() returned error: %v", err)
		}

		want := strings.Join(
			[]string{
				"az deployment sub",
				string(op),
				"--name",
				deployment,
				"--subscription",
				subId,
				"--location",
				location,
				"--template-file",
				filepath.Join(rr.Get(), bicep),
				"--out",
				"yaml",
				"--parameters",
				"@" + filepath.Join(rr.Get(), json),
			},
			" ")

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	}
}

func TestNewManagementGroup(t *testing.T) {
	t.Parallel()

	for _, op := range do.Operations() {
		got, err := deploy.New(dl.ManagementGroup, op, context, mgId, location, bicep, json)
		if err != nil {
			t.Errorf("New() returned error: %v", err)
		}

		want := strings.Join(
			[]string{
				"az deployment mg",
				string(op),
				"--name",
				deployment,
				"--management-group-id",
				mgId,
				"--location",
				location,
				"--template-file",
				filepath.Join(rr.Get(), bicep),
				"--out",
				"yaml",
				"--parameters",
				"@" + filepath.Join(rr.Get(), json),
			},
			" ")

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	}
}

//func TestResourceGroupAZCmdInvalidParameterFile(t *testing.T) {
//	t.Parallel()
//
//	rg, err := deploy.NewResourceGroup(subId, rgName, bicep, context)
//	if err != nil {
//		t.Errorf("NewResourceGroup() returned error: %v", err)
//	}
//	for _, op := range operation.Operations() {
//		if _, err := rg.AZCmd(op, "no param json"); err == nil {
//			t.Fatal("want error for invalid AZCmd, got nil")
//		}
//	}
//}
