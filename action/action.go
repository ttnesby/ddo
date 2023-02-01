package action

import (
	"context"
	"dagger.io/dagger"
	"ddo/alogger"
	"ddo/arg"
	"ddo/configuration"
	"ddo/path"
	"fmt"
	"github.com/tidwall/gjson"
)

var l = alogger.New()

type conctx struct {
	container *dagger.Container
	ctx       context.Context
}

type component struct {
	path   []string
	folder string
	tags   []string
}

func Do(specificationPath, actionsPath string, ctx context.Context) (e error) {

	l.Infof("Start dagger client")
	client, e := dagger.Connect(ctx)
	if e != nil {
		return l.Error(e)
	}

	defer func(client *dagger.Client) {
		_ = client.Close()
	}(client)

	dc, e := getContainer(client)
	if e != nil {
		return e
	}

	cont := conctx{
		container: dc,
		ctx:       ctx,
	}

	actionJson, e := getSpecification(specificationPath, cont)
	if e != nil {
		return e
	}

	l.Infof("Get actions : %s", actionsPath)
	actions := gjson.Get(actionJson, "actions."+actionsPath+"|@pretty")
	if !actions.Exists() {
		return l.Error(fmt.Errorf("no such path: %v", actionsPath))
	}

	if !gjson.Valid(actions.String()) {
		return l.Error(fmt.Errorf("resulting json-actions from path is invalid"))
	}

	l.Debugf("Actions \n%v ", actions)

	l.Debugf("Parse actions")
	components := parse(arg.LastActions(), actions)

	l.Debugf("Components: %v", components)

	doComponents(components, cont)

	l.Infof("Done!")

	return nil
}

func doComponents(components []component, c conctx) {
	for _, component := range components {
		yaml, err := c.exec(configuration.New(component.folder, component.tags).AsYaml())
		if err != nil {
			l.Errorf("Component %v with failing config \n%v", component.path, err)
		} else {
			l.Infof("Component %v \n%v", component.path, yaml)
		}
	}
}

func parse(lastActions []string, json gjson.Result) (components []component) {

	tmp := make([]component, 0)

	var loop func(lastActions []string, json gjson.Result)
	loop = func(lastActions []string, json gjson.Result) {
		if json.IsObject() {
			if objectIsComponent(json.Value().(map[string]interface{})) {
				tmp = append(
					tmp,
					createComponent(
						lastActions,
						json.Get("folder").String(),
						getTags(json.Get("tags").Array()),
					))
			} else {
				json.ForEach(func(key, value gjson.Result) bool {
					l.Debugf("key: %v", key)
					loop(append(lastActions, key.String()), value)
					return true
				})
			}
		}
	}

	loop(lastActions, json)

	return tmp
}

func createComponent(path []string, folder string, tags []string) component {
	return component{path: path, folder: folder, tags: tags}
}

func getTags(r []gjson.Result) (tags []string) {
	for _, v := range r {
		tags = append(tags, v.String())
	}
	return tags
}

func objectIsComponent(v map[string]interface{}) bool {
	_, hasFolder := v["folder"]
	_, hasTags := v["tags"]
	return hasFolder && hasTags
}

func getSpecification(specificationPath string, c conctx) (actionJson string, e error) {

	l.Infof("Reading action specification %v", specificationPath)
	actionJson, e = c.exec(configuration.New(specificationPath, nil).AsJson())
	if e != nil {
		return "", l.Error(e)
	}

	return actionJson, nil
}

func (c conctx) exec(cmd []string) (stdout string, e error) {
	return c.container.WithExec(cmd).Stdout(c.ctx)
}

func getContainer(client *dagger.Client) (*dagger.Container, error) {

	const (
		containerRef      = "docker.io/ttnesby/azbicue:latest"
		containerDotAzure = "/root/.azure"
		containerRepoRoot = "/rr"
	)

	l.Debugf("Verify and connect to host repository %s\n", path.RepoRoot())
	l.Debugf("Verify and connect to host %s\n", getDotAzurePath())
	if !path.AbsExists(getDotAzurePath()) {
		return nil, l.Error(fmt.Errorf("folder %s does not exist", getDotAzurePath()))
	}

	l.Debugf("Start container %s mounting [repo root, .azure]\n", containerRef)

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
