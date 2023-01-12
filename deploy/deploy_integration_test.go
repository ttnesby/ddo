//go:build integration

package deploy_test

import (
	"ddo/deploy"
	dl "ddo/deploy/level"
	do "ddo/deploy/operation"
	yaml "gopkg.in/yaml.v3"
	"os/exec"
	"strings"
	"testing"
)

const (
	itLocation = "norwayEast"
	itSubId    = "ca1e4592-6211-4c03-aac5-e681c1d1ea0d"
	itBicep    = "./test/data/resourceGroup.bicep"
	itJson     = "./test/data/resourceGroup.params.json"
)

func TestDeploymentSubscriptionValidate(t *testing.T) {
	//t.Parallel()

	//for _, op := range do.Operations() {
	azCmd, err := deploy.New(dl.Subscription, do.Validate, "rg-ddo-test", itSubId, itLocation, itBicep, itJson)
	if err != nil {
		t.Errorf("New() returned error: %v\n", err)
	}

	cmd := exec.Command(azCmd[0], azCmd[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cmd.CombinedOutput() for %s\nfailed with %s\n", azCmd, err)
	}

	data := make(map[string]interface{})
	if err = yaml.Unmarshal(out, &data); err != nil {
		t.Fatalf("yaml.Unmarshal returned error %v\n", err)
	}

	if provisioningState := data["properties"].(map[string]interface{})["provisioningState"]; provisioningState != "Succeeded" {
		t.Errorf("Deployment failed with provisioningState: %v", provisioningState)
	}
}

func TestDeploymentSubscriptionWhatIf(t *testing.T) {
	//t.Parallel()

	//for _, op := range do.Operations() {
	azCmd, err := deploy.New(dl.Subscription, do.WhatIf, "rg-ddo-test", itSubId, itLocation, itBicep, itJson)
	if err != nil {
		t.Errorf("New() returned error: %v\n", err)
	}

	//whatIf fails when started as a subprocess, so we have to run it as a shell command
	cmd := exec.Command("/bin/sh", "-c", strings.Join(azCmd, " "))
	if _, err = cmd.CombinedOutput(); err != nil {
		t.Fatalf("cmd.CombinedOutput() /bin/sh -c for %s\nfailed with %s\n", azCmd, err)
	}
}

func TestDeploymentSubscriptionDeploy(t *testing.T) {
	//t.Parallel()

	//for _, op := range do.Operations() {
	azCmd, err := deploy.New(dl.Subscription, do.Deploy, "rg-ddo-test", itSubId, itLocation, itBicep, itJson)
	if err != nil {
		t.Errorf("New() returned error: %v\n", err)
	}

	cmd := exec.Command(azCmd[0], azCmd[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cmd.CombinedOutput() for %s\nfailed with %s\n", azCmd, err)
	}

	data := make(map[string]interface{})
	if err = yaml.Unmarshal(out, &data); err != nil {
		t.Fatalf("yaml.Unmarshal failed with #{err}\n")
	}

	if provisioningState := data["properties"].(map[string]interface{})["provisioningState"]; provisioningState != "Succeeded" {
		t.Errorf("Deployment failed with provisioningState: %v", provisioningState)
	}
}
