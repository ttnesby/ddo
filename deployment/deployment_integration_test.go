//go:build integration

package deployment_test

import (
	dep "ddo/deployment"
	"testing"
)

const (
	itLocation = "norwayeast"
	itSubId    = "ca1e4592-6211-4c03-aac5-e681c1d1ea0d"
	itBicep    = "./test/data/resourceGroup.bicep"
	itJson     = "./test/data/resourceGroup.params.json"
)

func TestDeploymentSubscriptionValidate(t *testing.T) {
	t.Parallel()

	azCmd, err := dep.Validate(itBicep, itJson, dep.Subscription(itSubId, itLocation))
	if err != nil {
		t.Error(err)
	}

	if _, _, err = azCmd.Run(); err != nil {
		t.Error(err)
	}
}

// issue - something is wrong using exec - needs to be fixed
func TestDeploymentSubscriptionWhatIf(t *testing.T) {
	//t.Parallel()

	azCmd, err := dep.WhatIf(itBicep, itJson, dep.Subscription(itSubId, itLocation))
	if err != nil {
		t.Error(err)
	}

	if _, _, err = azCmd.Run(); err != nil {
		t.Error(err)
	}
}

func TestDeploymentSubscriptionDeploy(t *testing.T) {
	t.Parallel()

	azCmd, err := dep.Deploy(itBicep, itJson, dep.Subscription(itSubId, itLocation))
	if err != nil {
		t.Error(err)
	}

	if _, _, err = azCmd.Run(); err != nil {
		t.Error(err)
	}
}
