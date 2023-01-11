package deploylevel

import (
	"ddo/deployoperation"
	"ddo/reporoot"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func name(context string) string {
	sha1 := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(context))
	wo := strings.Split(sha1.String(), "-")
	return context + "-" + wo[0] + wo[4]
}

type ResourceGroup struct {
	Deployment        string
	SubscriptionId    uuid.UUID
	ResourceGroupName string
	Template          string
}

func NewResourceGroup(subId string, rgName string, templateFile string, context string) (ResourceGroup, error) {

	id, err := uuid.Parse(subId)
	if err != nil {
		return ResourceGroup{}, err
	}

	if _, err := os.Stat(filepath.Join(reporoot.Get(), templateFile)); err != nil {
		return ResourceGroup{}, err
	}

	return ResourceGroup{
		Deployment:        name(context),
		SubscriptionId:    id,
		ResourceGroupName: rgName,
		Template:          filepath.Join(reporoot.Get(), templateFile),
	}, nil
}

func (rg ResourceGroup) AZCmd(op deployoperation.Operation, parameterFile string) (string, error) {

	if _, err := os.Stat(filepath.Join(reporoot.Get(), parameterFile)); err != nil {
		return "", err
	}

	return strings.Join(
		[]string{
			"az deployment group",
			string(op),
			"--name",
			rg.Deployment,
			"--subscription",
			rg.SubscriptionId.String(),
			"--resource-group",
			rg.ResourceGroupName,
			"--template-file",
			rg.Template,
			"--out",
			"yaml",
			"--parameters",
			"@" + filepath.Join(reporoot.Get(), parameterFile),
		},
		" "), nil
}

type Subscription struct {
	Deployment     string
	SubscriptionId uuid.UUID
	Location       string
	Template       os.FileInfo
}

func NewSubscription(subId string, location string, templateFile string, context string) (Subscription, error) {

	id, err := uuid.Parse(subId)

	if err != nil {
		return Subscription{}, err
	}

	tf, err := os.Stat(templateFile)

	if err != nil {
		return Subscription{}, err
	}

	return Subscription{
		Deployment:     name(context),
		SubscriptionId: id,
		Location:       location,
		Template:       tf,
	}, nil
}

// func (sub Subscription) azCmd(c deployoperation.Category, parameterFile string) (string, error) {
// 	op, err := c.ToOperation()

// 	if err != nil {
// 		return "", err
// 	}

// 	pf, err := os.Stat(parameterFile)

// 	if err != nil {
// 		return "", err
// 	}

// 	return strings.Join(
// 		[]string{
// 			"az deployment sub",
// 			op,
// 			"--name",
// 			sub.Deployment,
// 			"--subscription",
// 			sub.SubscriptionId.String(),
// 			"--location",
// 			sub.Location,
// 			"--template-file",
// 			sub.Template.Name(),
// 			"--out",
// 			"yaml",
// 			"--parameters",
// 			"@" + pf.Name(),
// 		},
// 		" "), nil
// }

type ManagementGroup struct {
	Deployment        string
	ManagementGroupId uuid.UUID
	Location          string
	Template          os.FileInfo
}

func NewManagementGroup(mgId string, location string, templateFile string, context string) (ManagementGroup, error) {

	id, err := uuid.Parse(mgId)

	if err != nil {
		return ManagementGroup{}, err
	}

	tf, err := os.Stat(templateFile)

	if err != nil {
		return ManagementGroup{}, err
	}

	return ManagementGroup{
		Deployment:        name(context),
		ManagementGroupId: id,
		Location:          location,
		Template:          tf,
	}, nil
}

// func (mg ManagementGroup) azCmd(c deployoperation.Category, parameterFile string) (string, error) {
// 	op, err := c.ToOperation()

// 	if err != nil {
// 		return "", err
// 	}

// 	pf, err := os.Stat(parameterFile)

// 	if err != nil {
// 		return "", err
// 	}

// 	return strings.Join(
// 		[]string{
// 			"az deployment mg",
// 			op,
// 			"--name",
// 			mg.Deployment,
// 			"--management-group-id",
// 			mg.ManagementGroupId.String(),
// 			"--location",
// 			mg.Location,
// 			"--template-file",
// 			mg.Template.Name(),
// 			"--out",
// 			"yaml",
// 			"--parameters",
// 			"@" + pf.Name(),
// 		},
// 		" "), nil
// }
