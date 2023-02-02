package action

import (
	"context"
	"dagger.io/dagger"
	"ddo/alogger"
	"ddo/arg"
	"ddo/configuration"
	dep "ddo/deployment"
	"ddo/path"
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
	"sync"
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

func Do(ctx context.Context) (e error) {

	l.Infof("Start dagger client")
	//TODO the following should be enabled when debug mode
	//client, e := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
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

	actionSpec := path.ActionSpecification()
	l.Infof("Find ddo.cue %v", actionSpec)
	if len(actionSpec) == 0 || len(actionSpec) > 1 {
		return l.Error(fmt.Errorf("%d ddo.cue found", len(actionSpec)))
	}

	actionJson, e := getSpecification(actionSpec[0], cont)
	if e != nil {
		return e
	}

	l.Infof("Get actions: %s", arg.ActionsPath())
	actions := gjson.Get(actionJson, "actions."+arg.ActionsPath()+"|@pretty")
	if !actions.Exists() {
		return l.Error(fmt.Errorf("no such path: %v", arg.ActionsPath()))
	}

	if !gjson.Valid(actions.String()) {
		return l.Error(fmt.Errorf("resulting json-actions from path is invalid"))
	}

	l.Debugf("Parse actions")
	components := parse(arg.LastActions(), actions)
	l.Debugf("Components: %v", components)
	doComponents(arg.Operation(), components, cont)

	l.Infof("Done!")

	return nil
}

func configExport(component component, c conctx, wg *sync.WaitGroup) {
	defer wg.Done()

	yaml, err := c.exec(configuration.New(component.folder, component.tags).AsYaml())
	if err != nil {
		l.Errorf("component %v with failing config \n%v", component.path, err)
	} else {
		l.Infof("Component %v \n%v", component.path, yaml)
	}
}

//func showTimer(component string, stopTimer <-chan bool, wg *sync.WaitGroup) {
//
//	defer wg.Done()
//
//	ticker := time.Tick(time.Second)
//	start := time.Now()
//
//	stopSignal := false
//
//	go func() {
//		<-stopTimer
//		stopSignal = true
//	}()
//
//	fmt.Println()
//	for {
//		<-ticker
//		fmt.Printf("\r[%s]: %s", component, time.Now().Sub(start))
//		if stopSignal {
//			break
//		}
//	}
//	fmt.Println()
//}

//func configValidate(component component, c conctx, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	templatePath, toJsonCmd, tmpFile, destination, err := componentDetails(component, c)
//	if err != nil {
//		return
//	}
//
//	c.container = c.container.WithExec(toJsonCmd) // need the updated container with tmp file for az cli cmd
//
//	azCmd, err := dep.Validate(templatePath, tmpFile, destination)
//	if err != nil {
//		l.Errorf("could not create az cli validation command %v", err)
//		return
//	}
//
//	yaml, err := c.exec(azCmd)
//	if err != nil {
//		l.Errorf("%v failed\n%v", component.path, err)
//	} else {
//		l.Infof("%v done", component.path)
//		_ = yaml
//		//l.Infof("Component %v \n%v", component.path, yaml)
//	}
//}

func configAzCmd(
	component component,
	getAzCmd func(string, string, dep.ADestination) (dep.AzCli, error),
	c conctx,
	wg *sync.WaitGroup) {

	defer wg.Done()

	templatePath, toJsonCmd, tmpFile, destination, err := componentDetails(component, c)
	if err != nil {
		return
	}

	c.container = c.container.WithExec(toJsonCmd) // need the updated container with tmp file for az cli cmd

	azCmd, err := getAzCmd(templatePath, tmpFile, destination)
	if err != nil {
		l.Errorf("%v failed\n%v", component.path, err)
		return
	}

	yaml, err := c.exec(azCmd)
	if err != nil {
		l.Errorf("%v failed\n%v", component.path, err)
	} else {
		//l.Infof("%v done", component.path)
		l.Infof("%v \n%v", component.path, yaml)
	}
}

//func configWhatIf(component component, c conctx, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	templatePath, toJsonCmd, tmpFile, destination, err := componentDetails(component, c)
//	if err != nil {
//		return
//	}
//
//	c.container = c.container.WithExec(toJsonCmd) // need the updated container with tmp file for az cli cmd
//
//	azCmd, err := dep.WhatIf(templatePath, tmpFile, destination)
//	if err != nil {
//		l.Errorf("could not create az cli validation command %v", err)
//		return
//	}
//
//	yaml, err := c.exec(azCmd)
//	if err != nil {
//		l.Errorf("Component %v with failing validation \n%v", component.path, err)
//	} else {
//		l.Infof("Component %v \n%v", component.path, yaml)
//	}
//}

//func configDeploy(component component, c conctx, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	templatePath, toJsonCmd, tmpFile, destination, err := componentDetails(component, c)
//	if err != nil {
//		return
//	}
//
//	c.container = c.container.WithExec(toJsonCmd) // need the updated container with tmp file for az cli cmd
//
//	azCmd, err := dep.Deploy(templatePath, tmpFile, destination)
//	if err != nil {
//		l.Errorf("could not create az cli validation command %v", err)
//		return
//	}
//
//	yaml, err := c.exec(azCmd)
//	if err != nil {
//		l.Errorf("Component %v with failing validation \n%v", component.path, err)
//	} else {
//		l.Infof("%v \n%v", component.path, yaml)
//	}
//}

func componentDetails(component component, c conctx) (
	templatePath string,
	toJsonCmd configuration.CueCli,
	tmpFile string,
	destination dep.ADestination,
	e error) {

	anError := func(err error) (string, configuration.CueCli, string, dep.ADestination, error) {
		return "", configuration.CueCli{}, "", nil, l.Error(err)
	}

	config := configuration.New(component.folder, component.tags)

	// template path
	templatePath, err := c.exec(config.ElementsAsText([]string{"templatePath"}))
	if err != nil {
		return anError(fmt.Errorf("%v failed\n%v", component.path, err))
	}
	l.Debugf("%v templatePath %v", component.path, templatePath)

	// target
	targetJson, err := c.exec(config.ElementsAsJson([]string{"target"}))
	if err != nil {
		return anError(fmt.Errorf("could not extract target %v", err))
	}
	l.Debugf("%v targetJson %v", component.path, targetJson)

	// parameters as json file
	tmpFile = path.ContainerTmpJson()
	toJsonCmd = config.ElementsToTmpJsonFile(tmpFile, []string{"parameters"})

	return templatePath, toJsonCmd, tmpFile, resolveTarget(targetJson), nil
}

func resolveTarget(json string) dep.ADestination {
	if gjson.Get(json, "resourceGroup").Exists() {
		return dep.ResourceGroup(
			gjson.Get(json, "resourceGroup.name").String(),
			gjson.Get(json, "resourceGroup.inSubscriptionId").String(),
		)
	} else if gjson.Get(json, "subscription").Exists() {
		return dep.Subscription(
			gjson.Get(json, "subscription.id").String(),
			gjson.Get(json, "subscription.location").String(),
		)
	} else {
		return dep.ManagementGroup(
			gjson.Get(json, "managementGroup.id").String(),
			gjson.Get(json, "managementGroup.location").String(),
		)
	}
}

func doComponents(operation string, components []component, c conctx) {

	var cwg sync.WaitGroup

	//var twg sync.WaitGroup
	//stopTimer := make(chan bool)
	//twg.Add(1)
	//go showTimer("Processing components", stopTimer, &twg)
	validate := func(template, parameters string, dst dep.ADestination) (dep.AzCli, error) {
		return dep.Validate(template, parameters, dst)
	}

	whatif := func(template, parameters string, dst dep.ADestination) (dep.AzCli, error) {
		return dep.WhatIf(template, parameters, dst)
	}

	deploy := func(template, parameters string, dst dep.ADestination) (dep.AzCli, error) {
		return dep.Deploy(template, parameters, dst)
	}

	for _, component := range components {
		switch operation {
		case "ce":
			cwg.Add(1)
			go configExport(component, c, &cwg)
		case "va":
			cwg.Add(1)
			go configAzCmd(component, validate, c, &cwg)
		case "if":
			cwg.Add(1)
			go configAzCmd(component, whatif, c, &cwg)
		case "de":
			cwg.Add(1)
			//TODO Need to manage dependencies between components
			go configAzCmd(component, deploy, c, &cwg)
		default:
			l.Errorf("Operation %v not supported", operation)
		}
	}

	cwg.Wait()
	//stopTimer <- true
	//twg.Wait()
	//close(stopTimer)
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
	r, err := c.container.WithExec(cmd).Stdout(c.ctx)
	return strings.TrimRight(r, "\r\n"), err
}

func getContainer(client *dagger.Client) (*dagger.Container, error) {

	const (
		containerRef      = "docker.io/ttnesby/azbicue:latest"
		containerDotAzure = "/root/.azure"
		containerRepoRoot = "/rr"
	)

	l.Debugf("Verify and connect to host repository %s", path.RepoRoot())
	l.Debugf("Verify and connect to host %s", getDotAzurePath())
	if !path.AbsExists(getDotAzurePath()) {
		return nil, l.Error(fmt.Errorf("folder %s does not exist", getDotAzurePath()))
	}

	l.Debugf("Start container %s mounting [repo root, .azure]", containerRef)

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
