//go:build unit

package destination_test

import (
	"ddo/deployment/destination"
	"github.com/google/go-cmp/cmp"
	"testing"
)

const (
	location  = "norwayeast"
	validId   = "6ba7b810-9dad-11d1-80b4-00c04fd43022"
	invalidId = "6ba7"
	rgName    = "myRg"
)

func TestResourceGroup(t *testing.T) {
	t.Parallel()

	wantDestination := "group"
	wantParams := []string{
		"--subscription",
		validId,
		"--resource-group",
		rgName,
	}
	got, err := destination.ResourceGroup(rgName, validId)

	if err != nil {
		t.Errorf("ResourceGroup() returned error: %v", err)
	}

	gotDestination, gotParams := got.AzCli()

	if gotDestination != wantDestination {
		t.Errorf("ResourceGroup() returned destination %v, want %v", gotDestination, wantDestination)
	}

	if !cmp.Equal(wantParams, gotParams) {
		t.Error(cmp.Diff(wantParams, gotParams))
	}
}

func TestResourceGroupInvalidId(t *testing.T) {
	t.Parallel()

	_, err := destination.ResourceGroup(rgName, invalidId)

	if err == nil {
		t.Errorf("ResourceGroup() should return error when invalid id, got nil")
	}
}

func TestSubscription(t *testing.T) {
	t.Parallel()

	wantDestination := "sub"
	wantParams := []string{
		"--subscription",
		validId,
		"--location",
		location,
	}
	got, err := destination.Subscription(validId, location)

	if err != nil {
		t.Errorf("Subscription() returned error: %v", err)
	}

	gotDestination, gotParams := got.AzCli()

	if gotDestination != wantDestination {
		t.Errorf("Subscription() returned destination %v, want %v", gotDestination, wantDestination)
	}

	if !cmp.Equal(wantParams, gotParams) {
		t.Error(cmp.Diff(wantParams, gotParams))
	}
}

func TestManagementGroup(t *testing.T) {
	t.Parallel()

	wantDestination := "mg"
	wantParams := []string{
		"--management-group-id",
		validId,
		"--location",
		location,
	}
	got, err := destination.ManagementGroup(validId, location)

	if err != nil {
		t.Errorf("ManagementGroup() returned error: %v", err)
	}

	gotDestination, gotParams := got.AzCli()

	if gotDestination != wantDestination {
		t.Errorf("ManagementGroup() returned destination %v, want %v", gotDestination, wantDestination)
	}

	if !cmp.Equal(wantParams, gotParams) {
		t.Error(cmp.Diff(wantParams, gotParams))
	}
}
