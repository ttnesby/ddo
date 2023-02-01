package deployment_test

import (
	"ddo/alogger"
	dep "ddo/deployment"
	"ddo/path"
	"ddo/util"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

const (
	location      = "norwayeast"
	mgId          = "6ba7b810-9dad-11d1-80b4-00c04fd43022"
	subId         = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	rgName        = "myRg"
	validBicep    = "./test/data/resourceGroup.bicep"
	validParams   = "./test/data/resourceGroup.params.json"
	invalidUUID   = "6ba7"
	invalidBicep  = "./test/n/a/invalid.bicep"
	invalidParams = "./test/n/a/invalid.params.json"
)

func TestMain(m *testing.M) {
	//setup
	alogger.Disable()
	code := m.Run()
	//shutdown
	os.Exit(code)
}

func destAndOps(templatePath, paramsPath string, aDest dep.ADestination) error {

	theDest, err := aDest()
	if err != nil {
		return err
	}
	dest, destParams := theDest.AzCli()

	ops := []func() (dep.AzCli, error){
		func() (dep.AzCli, error) { return dep.Validate(templatePath, paramsPath, aDest) },
		func() (dep.AzCli, error) { return dep.WhatIf(templatePath, paramsPath, aDest) },
		func() (dep.AzCli, error) { return dep.Deploy(templatePath, paramsPath, aDest) },
	}

	tfp := path.RepoAbs(templatePath)
	if !path.AbsExists(tfp) {
		return fmt.Errorf("template file %s does not exist", tfp)
	}

	pfp := path.RepoAbs(paramsPath)
	if !path.AbsExists(pfp) {
		return fmt.Errorf("params file %s does not exist", pfp)
	}

	for _, op := range ops {
		got, err := op()
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		// check for what-if handling
		want := func() dep.AzCli {
			if got.IsWhatIf() {
				return []string{
					"az",
					"deployment",
					dest,
					got[3],
					got[4], // adding --what-if
					"--name",
					got[6],
					destParams[0],
					destParams[1],
					destParams[2],
					destParams[3],
					"--template-file",
					tfp,
					"--parameters",
					"@" + pfp,
					"--out",
					"yaml",
				}
			} else {
				return []string{
					"az",
					"deployment",
					dest,
					got[3],
					"--name",
					got[5],
					destParams[0],
					destParams[1],
					destParams[2],
					destParams[3],
					"--template-file",
					tfp,
					"--parameters",
					"@" + pfp,
					"--out",
					"yaml",
				}
			}
		}()

		if !cmp.Equal(want, got) {
			return fmt.Errorf("%v", cmp.Diff(want, got))
		}
	}
	return nil
}

func TestDeploymentToResourceGroup(t *testing.T) {
	t.Parallel()

	if err := destAndOps(validBicep, validParams, dep.ResourceGroup(rgName, subId)); err != nil {
		t.Errorf("DeploymentToResourceGroup() failed: %v", err)
	}
}

func TestDeploymentToResourceGroupInvalidUUID(t *testing.T) {
	t.Parallel()

	if err := destAndOps(validBicep, validParams, dep.ResourceGroup(rgName, invalidUUID)); err == nil {
		t.Error("want error for invalid UUID, got nil")
	}
}

func TestDeploymentToResourceGroupInvalidBicep(t *testing.T) {
	t.Parallel()

	if err := destAndOps(invalidBicep, validParams, dep.ResourceGroup(rgName, subId)); err == nil {
		t.Error("want error for invalid template, got nil")
	}
}

func TestDeploymentToResourceGroupInvalidParams(t *testing.T) {
	t.Parallel()

	if err := destAndOps(validBicep, invalidParams, dep.ResourceGroup(rgName, subId)); err == nil {
		t.Error("want error for invalid params, got nil")
	}
}

func TestDeploymentToSubscription(t *testing.T) {
	t.Parallel()

	if err := destAndOps(validBicep, validParams, dep.Subscription(subId, location)); err != nil {
		t.Errorf("DeploymentToSubscription() failed: %v", err)
	}
}

func TestDeploymentToManagementGroup(t *testing.T) {
	t.Parallel()

	if err := destAndOps(validBicep, validParams, dep.ManagementGroup(mgId, location)); err != nil {
		t.Errorf("DeploymentToManagementGroup() failed: %v", err)
	}
}

// integration tests

const (
	itLocation = "norwayeast"
	itSubId    = "ca1e4592-6211-4c03-aac5-e681c1d1ea0d"
	itBicep    = "./test/data/resourceGroup.bicep"
	itJson     = "./test/data/resourceGroup.params.json"
)

func TestDeploymentSubscriptionValidate(t *testing.T) {
	t.Parallel()

	util.SkipIfNotLoggedIntoNAVUTV(t)

	azCmd, err := dep.Validate(itBicep, itJson, dep.Subscription(itSubId, itLocation))
	if err != nil {
		t.Error(err)
	}

	if _, err = azCmd.Run(); err != nil {
		t.Error(err)
	}
}

// issue - something is wrong using exec - needs to be fixed
func TestDeploymentSubscriptionWhatIf(t *testing.T) {
	t.Parallel()

	util.SkipIfNotLoggedIntoNAVUTV(t)

	azCmd, err := dep.WhatIf(itBicep, itJson, dep.Subscription(itSubId, itLocation))
	if err != nil {
		t.Error(err)
	}

	if _, err = azCmd.Run(); err != nil {
		t.Error(err)
	}
}

func TestDeploymentSubscriptionDeploy(t *testing.T) {
	t.Parallel()

	util.SkipIfNotLoggedIntoNAVUTV(t)

	azCmd, err := dep.Deploy(itBicep, itJson, dep.Subscription(itSubId, itLocation))
	if err != nil {
		t.Error(err)
	}

	if _, err = azCmd.Run(); err != nil {
		t.Error(err)
	}
}