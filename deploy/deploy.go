package deploy

import (
	dl "ddo/deploy/level"
	do "ddo/deploy/operation"
	rr "ddo/reporoot"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type AzDeploy []string

func name(context string) string {
	sha1 := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(context))
	wo := strings.Split(sha1.String(), "-")
	return context + "-" + wo[0] + wo[4]
}

func verify(id, bicep, json string) (validId uuid.UUID, bicepPath string, jsonPath string, e error) {

	anError := func(e error) (uuid.UUID, string, string, error) {
		return uuid.Nil, "", "", e
	}

	i, err := uuid.Parse(id)
	if err != nil {
		return anError(err)
	}

	b := filepath.Join(rr.Get(), bicep)
	j := filepath.Join(rr.Get(), json)

	if _, err := os.Stat(b); err != nil {
		return anError(err)
	}
	if _, err := os.Stat(j); err != nil {
		return anError(err)
	}

	return i, b, j, nil
}

func New(
	level dl.Level,
	op do.Operation,
	context,
	id,
	rgOrLocation,
	templateFile,
	parameterFile string) (AzDeploy, error) {

	if !level.Valid() || !op.Valid() {
		return nil, fmt.Errorf("invalid level %s or operation %s", level, op)
	}

	i, b, j, err := verify(id, templateFile, parameterFile)
	if err != nil {
		return nil, err
	}

	prefix := []string{
		"az",
		"deployment",
		string(level),
		string(op),
		"--name",
		name(context),
	}

	infix := func() []string {
		switch level {
		case dl.ResourceGroup:
			return []string{
				"--subscription",
				i.String(),
				"--resource-group",
				rgOrLocation,
			}
		case dl.Subscription:
			return []string{
				"--subscription",
				i.String(),
				"--location",
				rgOrLocation,
			}
		case dl.ManagementGroup:
			return []string{
				"--management-group-id",
				i.String(),
				"--location",
				rgOrLocation,
			}
		}
		return nil // should never happen
	}()

	postfix := []string{
		"--template-file",
		b,
		"--out",
		"yaml",
		"--parameters",
		"@" + j,
	}

	return append(append(prefix, infix...), postfix...), nil
}

func (azCmd AzDeploy) Run() (asYaml map[string]interface{}, asByte []byte, e error) {

	isWhatIf := azCmd[3] == do.WhatIf.String()

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
