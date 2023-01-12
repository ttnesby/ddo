//go:build unit

package deploy_test

import (
	"ddo/deploy"
	dl "ddo/deploy/level"
	do "ddo/deploy/operation"
	rr "ddo/reporoot"
	"github.com/google/go-cmp/cmp"
	"path/filepath"
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

		want := []string{
			"az",
			"deployment",
			"group",
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
		}

		if !cmp.Equal(want, got) {
			t.Error(cmp.Diff(want, got))
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

		want := []string{
			"az",
			"deployment",
			"sub",
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
		}

		if !cmp.Equal(want, got) {
			t.Error(cmp.Diff(want, got))
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

		want := []string{
			"az",
			"deployment",
			"mg",
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
		}

		if !cmp.Equal(want, got) {
			t.Error(cmp.Diff(want, got))
		}
	}
}

func TestNewInvalidUUID(t *testing.T) {
	t.Parallel()

	for _, level := range dl.Levels() {
		for _, op := range do.Operations() {
			if _, err := deploy.New(level, op, context, "n/a", rgName, bicep, json); err == nil {
				t.Fatal("want error for invalid subId, got nil")
			}
		}
	}
}

func TestNewInvalidBicep(t *testing.T) {
	t.Parallel()

	for _, level := range dl.Levels() {
		for _, op := range do.Operations() {
			if _, err := deploy.New(level, op, context, mgId, rgName, "n/a", json); err == nil {
				t.Fatal("want error for invalid bicep, got nil")
			}
		}
	}
}

func TestNewInvalidJson(t *testing.T) {
	t.Parallel()

	for _, level := range dl.Levels() {
		for _, op := range do.Operations() {
			if _, err := deploy.New(level, op, context, mgId, rgName, bicep, "n/a"); err == nil {
				t.Fatal("want error for invalid json, got nil")
			}
		}
	}
}

func TestNewInvalidLevel(t *testing.T) {
	t.Parallel()

	for _, op := range do.Operations() {
		if _, err := deploy.New("n/a", op, context, subId, rgName, bicep, json); err == nil {
			t.Fatal("want error for invalid level, got nil")
		}
	}
}

func TestNewInvalidOperation(t *testing.T) {
	t.Parallel()

	for _, level := range dl.Levels() {
		if _, err := deploy.New(level, "n/a", context, subId, rgName, bicep, json); err == nil {
			t.Fatal("want error for invalid operation, got nil")
		}
	}
}
