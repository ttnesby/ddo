package action

import (
	"context"
	"dagger.io/dagger"
	"ddo/action/component"
	"ddo/alogger"
	"ddo/arg"
	"ddo/cuecli"
	"ddo/path"
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

func Do(ctx context.Context) (e error) {

	var client *dagger.Client

	l.Infof("start dagger client")
	if arg.DebugContainer() {
		client, e = dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout), dagger.WithWorkdir("."))
	} else {
		client, e = dagger.Connect(ctx, dagger.WithWorkdir("."))
	}
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

	selection, deployOrder, e := getSelectionAndDeployOrder(dc, ctx)
	if e != nil {
		return e
	}

	if e = doComponents(
		component.ActionsToComponents(
			selection,
			deployOrder,
			dc,
			ctx,
		),
	); e != nil {
		return e
	}

	l.Infof("done!")

	return nil
}

func getSelectionAndDeployOrder(
	container *dagger.Container,
	ctx context.Context) (selection, deployOrder gjson.Result, e error) {

	actionSpec := path.ActionSpecification()
	l.Infof("searched for ddo.cue %v", actionSpec)
	if len(actionSpec) == 0 || len(actionSpec) > 1 {
		return selection, deployOrder, l.Error(fmt.Errorf("%d ddo.cue file(s) found", len(actionSpec)))
	}

	l.Infof("reading action specification %v", actionSpec[0])
	actionJson, e := container.WithExec(
		cuecli.New(
			actionSpec[0],
			nil,
		).
			WithPackage("actions").
			AsJson(),
	).WithWorkdir(".").
		Stdout(ctx)

	actionJson = strings.TrimRight(actionJson, "\r\n")

	if e != nil {
		return selection, deployOrder, l.Error(e)
	}

	l.Infof("get selection: %v", arg.ActionsPath())
	selection = gjson.Get(actionJson, "actions."+arg.ActionsPath()+"|@pretty")
	if !selection.Exists() {
		return selection, deployOrder, l.Error(fmt.Errorf("no such path: %v", arg.ActionsPath()))
	}

	if !gjson.Valid(selection.String()) {
		return selection, deployOrder, l.Error(fmt.Errorf("resulting json-selection from path is invalid"))
	}

	l.Infof("get deployOrder")
	deployOrder = gjson.Get(actionJson, "deployOrder"+"|@pretty")
	if !deployOrder.Exists() {
		return selection, deployOrder, l.Error(fmt.Errorf("no deployOrder in ddo.cue"))
	}
	if !gjson.Valid(deployOrder.String()) {
		return selection, deployOrder, l.Error(fmt.Errorf("deployOrder is invalid"))
	}

	return selection, deployOrder, nil
}

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

func getContainer(client *dagger.Client) (*dagger.Container, error) {

	l.Debugf("connect to host repository [%s]", path.HostRepoRoot())
	l.Debugf("connect to host [%s]", path.HostDotAzure())
	l.Debugf("start container [%s] mounting [repo root, .azure]", path.ContainerRef)

	return client.Container().
		From(path.ContainerRef).
		WithMountedDirectory(path.ContainerRepoRoot, hostRepoRoot(client)).
		WithMountedDirectory(path.ContainerDotAzure, hostDotAzure(client)).
		WithWorkdir(path.ContainerRepoRoot), nil
}

func hostRepoRoot(c *dagger.Client) *dagger.Directory {
	return c.Host().Directory(path.HostRepoRoot(), dagger.HostDirectoryOpts{Exclude: []string{"build/"}})
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
