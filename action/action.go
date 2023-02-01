package action

import (
	"context"
	"dagger.io/dagger"
	"ddo/alogger"
	"ddo/configuration"
	"ddo/path"
	"fmt"
	"github.com/tidwall/gjson"
)

var l = alogger.New()

func Do(specificationPath, actionsPath string, ctx context.Context) (e error) {

	l.Infof("Start dagger client")
	client, e := dagger.Connect(ctx)
	if e != nil {
		return l.Error(e)
	}

	defer func(client *dagger.Client) {
		_ = client.Close()
	}(client)

	container, e := getContainer(client)
	if e != nil {
		return e
	}

	actionJson, e := getSpecification(specificationPath, container, ctx)
	if e != nil {
		return e
	}

	fullActionPath := "actions." + actionsPath
	l.Infof("Action path: %s", fullActionPath)
	actions := gjson.Get(actionJson, fullActionPath+"|@pretty")
	if !actions.Exists() {
		return l.Error(fmt.Errorf("no such path: %v", actionsPath))
	}

	l.Infof("actions: \n%v ", actions)

	l.Infof("Done!")

	return nil
}

func getSpecification(specificationPath string, container *dagger.Container, ctx context.Context) (actionJson string, e error) {

	l.Infof("Reading action specification %v", specificationPath)
	actionJson, e = container.WithExec(configuration.New(specificationPath, nil).AsJson()).Stdout(ctx)
	if e != nil {
		return "", l.Error(e)
	}

	return actionJson, nil
}

func getContainer(client *dagger.Client) (*dagger.Container, error) {

	const (
		containerRef      = "docker.io/ttnesby/azbicue:latest"
		containerDotAzure = "/root/.azure"
		containerRepoRoot = "/rr"
	)

	l.Infof("Get tooling container")
	l.Infof("Verify and connect to host repository %s\n", path.RepoRoot())
	l.Infof("Verify and connect to host %s\n", getDotAzurePath())
	if !path.AbsExists(getDotAzurePath()) {
		return nil, l.Error(fmt.Errorf("folder %s does not exist", getDotAzurePath()))
	}

	l.Infof("Start container %s mounting [repo root, .azure]\n", containerRef)

	return client.Container().
		From(containerRef).
		WithMountedDirectory(containerRepoRoot, hostRepoRoot(client)).
		WithMountedDirectory(containerDotAzure, hostDotAzure(client)).
		WithWorkdir(containerRepoRoot), nil
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
