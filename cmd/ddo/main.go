package main

import (
	"context"
	"dagger.io/dagger"
	"ddo/alogger"
	"ddo/configuration"
	"ddo/path"
	"fmt"
	"github.com/tidwall/gjson"
	"os"
	"strings"
)

var l = alogger.New()

func main() {
	exitCode := func() int {
		if err := build(context.Background()); err != nil {
			return 1
		}
		return 0
	}()
	os.Exit(exitCode)
}

func getDotAzurePath() string {
	return path.HomeAbs(".azure")
}

func hostRepoRoot(c *dagger.Client) *dagger.Directory {
	return c.Host().Directory(path.RepoRoot())
}

func hostDotAzure(c *dagger.Client) *dagger.Directory {
	return c.Host().Directory(getDotAzurePath(),
		dagger.HostDirectoryOpts{
			Include: []string{
				"azureProfile.json",
				"msal_http_cache.bin",
				"msal_token_cache.json",
				"service_principal_entries.json",
			},
		},
	)
}

func build(ctx context.Context) error {

	const (
		argProgramName = iota
		argPathToActionSpecification
		argActionsPath
		argMinNo = argActionsPath + 1
	)
	const (
		containerRef      = "docker.io/ttnesby/azbicue:latest"
		containerDotAzure = "/root/.azure"
		containerRepoRoot = "/rr"
	)

	l.Infof("Start program: %v", os.Args[argProgramName])

	if len(os.Args) < argMinNo {
		return l.Error(
			fmt.Errorf(
				"missing parameter(s) - usage: PROGRAM <path to action specification> <actions path>"))
	}

	l.Infof("Start dagger client")
	client, err := dagger.Connect(ctx)
	if err != nil {
		return l.Error(err)
	}

	defer func(client *dagger.Client) {
		_ = client.Close()
	}(client)

	l.Infof("Verify and connect to host repository %s\n", path.RepoRoot())
	l.Infof("Verify and connect to host %s\n", getDotAzurePath())
	if !path.AbsExists(getDotAzurePath()) {
		return l.Error(fmt.Errorf("folder %s does not exist", getDotAzurePath()))
	}

	l.Infof("Start container %s mounting [repo root, .azure]\n", containerRef)

	azbicue := client.Container().
		From(containerRef).
		WithMountedDirectory(containerRepoRoot, hostRepoRoot(client)).
		WithMountedDirectory(containerDotAzure, hostDotAzure(client)).
		WithWorkdir(containerRepoRoot)

	l.Infof("Reading action specification %v", os.Args[argPathToActionSpecification])
	actionsCmd := configuration.New(os.Args[argPathToActionSpecification], nil).AsJsonCmd()

	actionsJson, err := azbicue.WithExec(actionsCmd).Stdout(ctx)
	if err != nil {
		return l.Error(err)
	}

	actionsPath := "actions." + strings.Join(os.Args[argActionsPath:], ".")
	l.Infof("Action path: %s", actionsPath)
	actions := gjson.Get(actionsJson, actionsPath+"|@pretty")
	if !actions.Exists() {
		return l.Error(fmt.Errorf("no such path: %v", actionsPath))
	}

	l.Infof("actions: \n%v ", actions)

	l.Infof("Done!")

	return nil
}
