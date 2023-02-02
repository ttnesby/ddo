package deployment

import (
	"bytes"
	"ddo/alogger"
	"ddo/arg"
	de "ddo/deployment/destination"
	fp "ddo/path"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os/exec"
	"strings"
)

var l alogger.ALogger

func Init() {
	l = alogger.New(arg.InDebugMode())
}

type operation string

const (
	validate operation = "validate"
	whatIf   operation = "create --what-if"
	deploy   operation = "create"

	whatIfStart = 3
	whatIfEnd   = 5
)

type AzCli []string

func name(context string) string {
	sha1 := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(context))
	return sha1.String()
}

type ADestination func() (de.Destination, error)

func azDeploy(op operation, templatePath, parameterPath string, destination ADestination) (AzCli, error) {

	tfp := fp.RepoAbs(templatePath)
	pfp := fp.RepoAbs(parameterPath)

	theDest, err := destination()
	if err != nil {
		return nil, err
	}
	dest, destParams := theDest.AzCli()

	azCmd := []string{"az", "deployment", dest}
	azCmd = append(azCmd, strings.Split(string(op), " ")...) // due to whatIf
	azCmd = append(azCmd, "--name", name(tfp+pfp))
	azCmd = append(azCmd, destParams...)
	azCmd = append(azCmd, "--template-file", templatePath, "--parameters", "@"+parameterPath, "--out", "yaml")

	l.Debugf("azCmd: %v", azCmd)

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
	return strings.Join(azCmd[whatIfStart:whatIfEnd], " ") == string(whatIf)
}

func (azCmd AzCli) Run() (byte []byte, e error) {

	l.Debugf("azCmd: %v", azCmd)
	cmd := exec.Command(azCmd[0], azCmd[1:]...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(&stdoutBuf) //os.Stdout
	cmd.Stderr = io.MultiWriter(&stderrBuf) //os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, l.Error(fmt.Errorf("%v failed: %s", azCmd, err))
	}

	out, errStr := stdoutBuf.Bytes(), string(stderrBuf.Bytes())
	if errStr != "" {
		return nil, l.Error(fmt.Errorf("%v returned: %s", azCmd, errStr))
	}

	return out, nil
}
