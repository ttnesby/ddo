package component

import (
	"context"
	"dagger.io/dagger"
	"ddo/alogger"
	"ddo/arg"
	"ddo/configuration"
	dep "ddo/deployment"
	p "ddo/path"
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
	"sync"
)

type conctx struct {
	container *dagger.Container
	ctx       context.Context
}

type component struct {
	path   []string
	folder string
	tags   []string
	tech   conctx
}

var l alogger.ALogger

func Init() {
	l = alogger.New(arg.InDebugMode())
}

func (c component) config() configuration.CueCli {
	return configuration.New(c.folder, c.tags).WithPackage("deployment")
}

func (c component) templatePath() (path string, e error) {
	path, e = c.tech.exec(c.config().ElementsAsText([]string{"templatePath"}))
	switch {
	case e != nil:
		return "", l.Error(fmt.Errorf("%v failed\n%v", c.path, e))
	default:
		l.Debugf("%v templatePath %v", c.path, path)
		return path, nil
	}
}

func (c component) target() (dst dep.ADestination, e error) {
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

func (c component) paramsToTmpJsonFile() (path string, cmd configuration.CueCli) {
	path = p.ContainerTmpJson()
	cmd = c.config().ElementsToTmpJsonFile(path, []string{"parameters"})
	return path, cmd
}

func (c conctx) exec(cmd []string) (stdout string, e error) {
	r, err := c.container.WithExec(cmd).Stdout(c.ctx)
	return strings.TrimRight(r, "\r\n"), err
}

func (c component) exec(cmd []string, signalError chan<- bool) {
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

func (c component) configExport(signalError chan<- bool) {
	c.exec(c.config().AsYaml(), signalError)
}

func (c component) configDeploy(
	depOp func(string, string, dep.ADestination) (dep.AzCli, error),
	signalError chan<- bool) {

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

func (c component) do(signalError chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()

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
	case "ce":
		c.configExport(signalError)
	case "va":
		c.configDeploy(validate, signalError)
	case "if":
		c.configDeploy(whatif, signalError)
	case "de":
		c.configDeploy(deploy, signalError)
	default:
		signalError <- true
	}
}
