//go:build integration

package deploy_test

import (
	"ddo/deploy"
	dl "ddo/deploy/level"
	do "ddo/deploy/operation"
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

	azCmd, err := deploy.New(dl.Subscription, do.Validate, "rg-ddo-test", itSubId, itLocation, itBicep, itJson)
	if err != nil {
		t.Error(err)
	}

	if _, _, err = azCmd.Run(); err != nil {
		t.Error(err)
	}
}

func TestDeploymentSubscriptionWhatIf(t *testing.T) {
	//t.Parallel()

	azCmd, err := deploy.New(dl.Subscription, do.WhatIf, "rg-ddo-test", itSubId, itLocation, itBicep, itJson)
	if err != nil {
		t.Error(err)
	}

	if _, _, err = azCmd.Run(); err != nil {
		t.Error(err)
	}
}

func TestDeploymentSubscriptionDeploy(t *testing.T) {
	t.Parallel()

	azCmd, err := deploy.New(dl.Subscription, do.Deploy, "rg-ddo-test", itSubId, itLocation, itBicep, itJson)
	if err != nil {
		t.Error(err)
	}

	if _, _, err = azCmd.Run(); err != nil {
		t.Error(err)
	}
}
