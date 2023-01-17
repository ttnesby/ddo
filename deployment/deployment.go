package deployment

import (
	de "ddo/deployment/destination"
	fp "ddo/fullpath"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"os/exec"
	"strings"
)

type operation string

const (
	validate operation = "validate"
	whatIf   operation = "create --what-if"
	deploy   operation = "create"
)

type AzCli []string

func name(context string) string {
	sha1 := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(context))
	return sha1.String()
}

type ADestination func() (de.Destination, error)

func azDeploy(op operation, templatePath, parameterPath string, destination ADestination) (AzCli, error) {

	tfp, err := fp.Get(templatePath)
	if err != nil {
		return nil, err
	}

	pfp, err := fp.Get(parameterPath)
	if err != nil {
		return nil, err
	}

	theDest, err := destination()
	if err != nil {
		return nil, err
	}
	dest, destParams := theDest.AzCli()

	prefix := []string{"az", "deployment", dest, string(op), "--name", name(tfp + pfp)}
	postfix := []string{"--template-file", tfp, "--parameters", "@" + pfp, "--out", "yaml"}

	return append(append(prefix, destParams...), postfix...), nil
}

func ResourceGroup(name, inSubscriptionId string) ADestination {
	return func() (de.Destination, error) {
		return de.ResourceGroup(name, inSubscriptionId)
	}
}

func Subscription(id, location string) ADestination {
	return func() (de.Destination, error) {
		return de.Subscription(id, location)
	}
}

func ManagementGroup(id, location string) ADestination {
	return func() (de.Destination, error) {
		return de.ManagementGroup(id, location)
	}
}

func Validate(templatePath, parameterPath string, destination ADestination) (AzCli, error) {
	return azDeploy(validate, templatePath, parameterPath, destination)
}

func WhatIf(templatePath, parameterPath string, destination ADestination) (AzCli, error) {
	return azDeploy(whatIf, templatePath, parameterPath, destination)
}

func Deploy(templatePath, parameterPath string, destination ADestination) (AzCli, error) {
	return azDeploy(deploy, templatePath, parameterPath, destination)
}

func (azCmd AzCli) Run() (asYaml map[string]interface{}, asByte []byte, e error) {

	isWhatIf := operation(azCmd[3]) == whatIf

	cmd := func() *exec.Cmd {
		if isWhatIf {
			return exec.Command(azCmd[0], azCmd[1:]...)
		} else {
			return exec.Command("/bin/sh", "-c", strings.Join(azCmd, " "))
		}
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, nil, fmt.Errorf("Run() for %v\nreturned error %s\n", azCmd, err)
	}

	// only hard-to-parse output from what-if, leaving as is
	if isWhatIf {
		return nil, out, nil
	}

	// parse the output
	data := make(map[string]interface{})
	if err = yaml.Unmarshal(out, &data); err != nil {
		return nil, nil, fmt.Errorf("yaml.Unmarshal() of %v\nreturned error %v\n", out, err)
	}

	return data, nil, nil
}
