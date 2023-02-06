package component

import (
	"context"
	"dagger.io/dagger"
	"ddo/alogger"
	"ddo/arg"
	del "ddo/azcli/delete"
	dep "ddo/azcli/deployment"
	"ddo/configuration"
	p "ddo/path"
	"ddo/util"
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
	"sync"
)

type conctx struct {
	container *dagger.Container
	ctx       context.Context
}

type Component struct {
	path   []string
	folder string
	tags   []string
	tech   conctx
}

var l alogger.ALogger

func Init() {
	l = alogger.New(arg.InDebugMode())
}

func (c Component) config() configuration.CueCli {
	return configuration.New(c.folder, c.tags).WithPackage("deployment")
}

func (c Component) templatePath() (path string, e error) {
	path, e = c.tech.exec(c.config().ElementsAsText([]string{"templatePath"}))
	switch {
	case e != nil:
		return "", l.Error(fmt.Errorf("%v failed\n%v", c.path, e))
	default:
		l.Debugf("%v templatePath %v", c.path, path)
		return path, nil
	}
}

func (c Component) resourceId() (path string, e error) {
	path, e = c.tech.exec(c.config().ElementsAsText([]string{"#resourceId"}))
	switch {
	case e != nil:
		return "", l.Error(fmt.Errorf("%v failed\n%v", c.path, e))
	default:
		l.Debugf("%v #resourceId %v", c.path, path)
		return path, nil
	}
}

func (c Component) target() (dst dep.ADestination, e error) {
	targetJson, e := c.tech.exec(c.config().ElementsAsJson([]string{"target"}))
	switch {
	case e != nil:
		return nil, l.Error(fmt.Errorf("%v failed\n%v", c.path, e))
	default:
		l.Debugf("%v target %v", c.path, targetJson)
		switch {
		case gjson.Get(targetJson, "resourceGroup").Exists():
			dst = dep.ResourceGroup(
				gjson.Get(targetJson, "resourceGroup.name").String(),
				gjson.Get(targetJson, "resourceGroup.inSubscriptionId").String(),
			)
		case gjson.Get(targetJson, "subscription").Exists():
			dst = dep.Subscription(
				gjson.Get(targetJson, "subscription.id").String(),
				gjson.Get(targetJson, "subscription.location").String(),
			)
		default:
			dst = dep.ManagementGroup(
				gjson.Get(targetJson, "managementGroup.id").String(),
				gjson.Get(targetJson, "managementGroup.location").String(),
			)
		}
	}
	return dst, nil
}

func (c Component) paramsToTmpJsonFile() (path string, cmd configuration.CueCli) {
	path = p.ContainerTmpJson()
	cmd = c.config().ElementsToTmpJsonFile(path, []string{"parameters"})
	return path, cmd
}

// -- need to deal with remove and lack of feedback - no, fooled by cache...
func (c conctx) exec(cmd []string) (stdout string, e error) {
	r, err := c.container.WithExec(cmd).Stdout(c.ctx)
	return strings.TrimRight(r, "\r\n"), err
}

func (c Component) exec(cmd []string, signalError chan<- bool) {
	res, err := c.tech.exec(cmd)
	switch {
	case err != nil:
		signalError <- true
		l.Errorf("%v failed \n%v", c.path, err)
	case arg.NoResultDisplay():
		l.Infof("%v done", c.path)
	default:
		l.Infof("%v \n%v", c.path, res)
	}
}

func (c Component) remove(signalError chan<- bool) {

	l.Debugf("%v remove", c.path)

	//TODO - in case of type resourceGroups/ - do we need --force? In case of hosted resource cannot be deleted
	// and use of --force on the resource group will remediate the situation

	resourceId, err := c.resourceId()
	if err != nil {
		signalError <- true
		return
	}

	azCmd := del.ResourceId(resourceId)
	c.exec(azCmd, signalError)
}

func (c Component) configExport(signalError chan<- bool) {
	l.Debugf("%v configDeploy", c.path)
	c.exec(c.config().AsYaml(), signalError)
}

func (c Component) configDeploy(
	depOp func(string, string, dep.ADestination) (dep.AzCli, error),
	signalError chan<- bool) {

	l.Debugf("%v configDeploy", c.path)

	tmpJsonFile, exportCmd := c.paramsToTmpJsonFile()
	c.tech.container = c.tech.container.WithExec(exportCmd) // config params exported to tmp json file

	templatePath, err := c.templatePath()
	if err != nil {
		signalError <- true
		return
	}

	dst, err := c.target()
	if err != nil {
		signalError <- true
		return
	}

	azCmd, err := depOp(templatePath, tmpJsonFile, dst)
	if err != nil {
		signalError <- true
		return
	}

	c.exec(azCmd, signalError)
}

func (c Component) Do(signalError chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()

	l.Infof("%v %s", c.path, arg.Operation())

	validate := func(template, parameters string, dst dep.ADestination) (dep.AzCli, error) {
		return dep.Validate(template, parameters, dst)
	}

	whatif := func(template, parameters string, dst dep.ADestination) (dep.AzCli, error) {
		return dep.WhatIf(template, parameters, dst)
	}

	deploy := func(template, parameters string, dst dep.ADestination) (dep.AzCli, error) {
		return dep.Deploy(template, parameters, dst)
	}

	switch arg.Operation() {
	case arg.OpCE:
		c.configExport(signalError)
	case arg.OpVA:
		c.configDeploy(validate, signalError)
	case arg.OpIF:
		c.configDeploy(whatif, signalError)
	case arg.OpDE:
		c.configDeploy(deploy, signalError)
	case arg.OpRE:
		c.remove(signalError)
	default:
		signalError <- true
	}
}

func ActionsToComponents(
	selection, deployOrder gjson.Result,
	container *dagger.Container,
	ctx context.Context) (components [][]Component) {

	components = orderComponents(
		parse(
			selection,
			conctx{container: container, ctx: ctx},
		),
		deployOrder,
	)

	return components
}

func parse(json gjson.Result, cc conctx) (components []Component) {

	tmp := make([]Component, 0)
	lastActions := arg.LastActions()

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
						cc,
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

func orderComponents(components []Component, deployOrder gjson.Result) (ordCo [][]Component) {

	l.Debugf("group components %v", components)
	// make a list where each element is a list of components that can be deployed in parallel

	// order is not relevant for ce, va or if
	if op := arg.Operation(); op == arg.OpCE || op == arg.OpVA || op == arg.OpIF {
		l.Debugf("due operation %s, no order required", op)
		return append(ordCo, components)
	}

	for _, le := range deployOrder.Array() {
		var group []Component
		for _, co := range le.Array() {
			for _, c := range components {
				noOfPathElems := len(c.path)
				if c.path[noOfPathElems-1] == co.String() {
					group = append(group, c)
				}
			}
		}
		ordCo = append(ordCo, group)
	}

	// order must be reversed
	if op := arg.Operation(); op == arg.OpRE {
		l.Debugf("due operation %s, order is reversed", op)
		util.ReverseSlice(ordCo)
		return ordCo
	}

	l.Debugf("Ordered components %v", ordCo)
	return ordCo
}

func objectIsComponent(v map[string]interface{}) bool {
	_, hasFolder := v["folder"]
	_, hasTags := v["tags"]
	return hasFolder && hasTags
}

func createComponent(path []string, folder string, tags []string, cc conctx) Component {
	return Component{path: path, folder: folder, tags: tags, tech: cc}
}

func getTags(r []gjson.Result) (tags []string) {
	for _, v := range r {
		tags = append(tags, v.String())
	}
	return tags
}
