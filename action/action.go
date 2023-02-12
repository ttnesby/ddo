package action

import (
	"context"
	"dagger.io/dagger"
	"ddo/action/component"
	"ddo/alogger"
	"ddo/arg"
	"ddo/azcli/resource"
	"ddo/cuecli"
	"ddo/path"
	"encoding/base64"
	"fmt"
	"github.com/tidwall/gjson"
	"os"
	"strings"
	"sync"
	"time"
)

var l alogger.ALogger

func Init() {
	l = alogger.New(arg.InDebugMode())
}

// Do
// process the actions specification defined by ddo.cue
//
// * A recursive search from repo root will find a single ddo.cue
//
// * The chosen operation will be applied on each component selected by the action path
//
// * All data injections are resolved before operation invocation
//
// * In case of deploy, the ddo.cue-deployOrder defines the order of components
//
// * In case of evomer, a reversed deployOrder is applied
//
// * A group of components in deployOrder are invoked in parallel
func Do(ctx context.Context) (e error) {

	l.Infof("start dagger client")

	client, e := func(debug bool) (*dagger.Client, error) {
		if debug {
			return dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout), dagger.WithWorkdir("."))
		}
		return dagger.Connect(ctx, dagger.WithWorkdir("."))
	}(arg.DebugContainer())
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

	actionJson, e := getActionJson(dc, ctx)
	if e != nil {
		return e
	}

	selection, e := jsonSelect(actionJson, "actions."+arg.ActionsPath()+"|@pretty")
	if e != nil {
		return e
	}

	deployOrder, e := jsonSelect(actionJson, "deployOrder"+"|@pretty")
	if e != nil {
		return e
	}

	//selection, deployOrder, e := getSelectionAndDeployOrder(actionJson)
	//if e != nil {
	//	return e
	//}

	components := component.ActionsToComponents(selection, deployOrder, dc, ctx)

	if e = resolveDataInjection(actionJson, deployOrder, components, dc, ctx); e != nil {
		return e
	}

	if e = doComponents(components); e != nil {
		return e
	}

	l.Infof("done!")

	return nil
}

func resolveDataInjection(
	actionJson string,
	deployOrder gjson.Result,
	components [][]component.Component,
	container *dagger.Container,
	ctx context.Context) (e error) {

	getData := func(actionPath, dataPath string) (data string, e error) {

		fp := fmt.Sprintf("actions.%s.%s|@pretty", arg.OpCE, actionPath)
		selection, e := jsonSelect(actionJson, fp)
		if e != nil {
			return data, e
		}
		r := component.ActionsToComponents(selection, deployOrder, container, ctx)
		// can only do data lookup on a single component
		if len(r) != 1 || len(r[0]) != 1 {
			return data, l.Error(fmt.Errorf("data injection must be based on a single component"))
		}
		diCo := r[0][0]
		// get the resourceId
		rId, e := diCo.ResourceId()
		if e != nil {
			return data, e
		}

		//TODO need error handling and existence check
		json, _ := container.WithExec(resource.Show(rId)).Stdout(ctx)

		if len(json) == 0 {
			return "", nil
		}

		if dataPath == "b64" {
			return base64.StdEncoding.EncodeToString([]byte(json)), nil
		}

		data, _ = jsonSelectString(json, dataPath)

		return data, nil
	}

	l.Infof("resolve data injections for components")
	for _, group := range components {
		for _, co := range group {
			e = co.DataInjection(getData)
		}
	}
	return nil
}

func getActionJson(
	container *dagger.Container,
	ctx context.Context) (actionJson string, e error) {

	actionSpec := path.ActionSpecification()
	l.Infof("searched for ddo.cue %v", actionSpec)
	if len(actionSpec) == 0 || len(actionSpec) > 1 {
		return "", l.Error(fmt.Errorf("%d ddo.cue file(s) found", len(actionSpec)))
	}

	l.Infof("reading action specification %v", actionSpec[0])
	actionJson, e = container.WithExec(
		cuecli.New(
			actionSpec[0],
			nil,
		).
			WithPackage("actions").
			AsJson(),
	).WithWorkdir(".").
		Stdout(ctx)

	if e != nil {
		return "", l.Error(e)
	}

	return strings.TrimRight(actionJson, "\r\n"), nil
}

func jsonSelectString(actionJson, aPath string) (selection string, e error) {

	l.Infof("get selection: %v", aPath)
	sel := gjson.Get(actionJson, aPath)
	if !sel.Exists() {
		return selection, l.Error(fmt.Errorf("no such path: %v", aPath))
	}
	return sel.String(), nil
}

func jsonSelect(actionJson, aPath string) (selection gjson.Result, e error) {

	l.Infof("get selection: %v", aPath)
	selection = gjson.Get(actionJson, aPath)
	if !selection.Exists() {
		return selection, l.Error(fmt.Errorf("no such path: %v", aPath))
	}
	if !gjson.Valid(selection.String()) {
		return selection, l.Error(fmt.Errorf("resulting json-selection from path is invalid"))
	}
	return selection, nil
}

// doComponents
// iterates each group of components. Components in a group are invoked in parallel.
// Since a group is a prerequisite for later groups, the iteration is terminated in case of
// failure for operation `deployment` or `evomer`
func doComponents(groups [][]component.Component) error {

	l.Debugf("processing %v ", groups)
	var mu sync.Mutex
	noOfErrors := 0

	for _, group := range groups {

		var cwg sync.WaitGroup

		signalError := make(chan bool)
		stopListening := false

		go func() {
			for {
				<-signalError
				if stopListening {
					break
				}
				mu.Lock()
				noOfErrors++
				mu.Unlock()
			}
		}()

		for _, co := range group {
			cwg.Add(1)
			go co.Do(signalError, &cwg)
		}

		cwg.Wait()
		// TODO in case of evomer and non #resourceId - not catching then signal, due to the if-break?
		// with a little more wait, the error signal is picked up.It's ok with va, if - slower stuff
		time.Sleep(25 * time.Millisecond)
		stopListening = true
		close(signalError)

		if op := arg.Operation(); (op == arg.OpDE || op == arg.OpRE) && noOfErrors > 0 {
			// due to order dependency for these operations - stop now
			l.Infof("Error in prerequisites operation(s) - stopping")
			break
		}
	}

	if noOfErrors > 0 {
		return l.Error(fmt.Errorf("%v component(s) failed", noOfErrors))
	}
	return nil
}

// getContainer
// retrieves a container with the required tool set;
//
// - az cli
//
// - cue cli
//
// - bicep.
//
// Repo root and .azure (inherit host logins) are mounted
func getContainer(client *dagger.Client) (*dagger.Container, error) {

	l.Debugf("connect to host repository [%s]", path.HostRepoRoot())
	l.Debugf("connect to host [%s]", path.HostDotAzure())
	l.Debugf("start container [%s] mounting [repo root, .azure]", arg.ContainerRef())

	return client.Container().
		From(arg.ContainerRef()).
		WithMountedDirectory(path.ContainerRepoRoot, hostRepoRoot(client)).
		WithMountedDirectory(path.ContainerDotAzure, hostDotAzure(client)).
		WithWorkdir(path.ContainerRepoRoot), nil
}

func hostRepoRoot(c *dagger.Client) *dagger.Directory {
	return c.Host().Directory(path.HostRepoRoot(), dagger.HostDirectoryOpts{Exclude: []string{"output/"}})
}

func hostDotAzure(c *dagger.Client) *dagger.Directory {
	return c.Host().Directory(path.HostDotAzure(),
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
