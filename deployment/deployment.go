package deployment

import (
	"bytes"
	de "ddo/deployment/destination"
	fp "ddo/fullpath"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
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

	azCmd := []string{"az", "deployment", dest}
	azCmd = append(azCmd, strings.Split(string(op), " ")...) // due to whatIf
	azCmd = append(azCmd, "--name", name(tfp+pfp))
	azCmd = append(azCmd, destParams...)
	azCmd = append(azCmd, "--template-file", tfp, "--parameters", "@"+pfp, "--out", "yaml")

	return azCmd, nil
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

func (azCmd AzCli) IsWhatIf() bool {
	// !!observe, starting index : starting index + no of elements to get
	return strings.Join(azCmd[3:5], " ") == string(whatIf)
}

func (azCmd AzCli) Run() (byte []byte, e error) {

	cmd := exec.Command(azCmd[0], azCmd[1:]...)
	//out, err := cmd.CombinedOutput()
	//if err != nil {
	//	return nil, nil, fmt.Errorf("Run() for %v\nreturned error %s\n", azCmd, err)
	//}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("cmd.Run() of %v\nfailed with %s\n", azCmd, err)
	}

	out, errStr := stdoutBuf.Bytes(), string(stderrBuf.Bytes())
	//fmt.Printf("\nout:\n%s\nerr:\n%s\n", out, errStr)
	if errStr != "" {
		return nil, fmt.Errorf("%v\nreturned error %s\n", azCmd, errStr)
	}

	// parse the output
	//data := make(map[string]interface{})
	//if err = yaml.Unmarshal(out, &data); err != nil {
	//	return nil, nil, fmt.Errorf("yaml.Unmarshal() of %v\nreturned error %v\n", out, err)
	//}

	return out, nil
}
