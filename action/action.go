package action

import (
	"context"
	"dagger.io/dagger"
	"ddo/action/component"
	"ddo/alogger"
	"ddo/arg"
	"ddo/configuration"
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

	l.Infof("Start dagger client")
	if arg.DebugContainer() {
		client, e = dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	} else {
		client, e = dagger.Connect(ctx)
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

	l.Infof("Done!")

	return nil
}

func getSelectionAndDeployOrder(
	container *dagger.Container,
	ctx context.Context) (selection, deployOrder gjson.Result, e error) {

	actionSpec := path.ActionSpecification()
	l.Infof("Searched ddo.cue %v", actionSpec)
	if len(actionSpec) == 0 || len(actionSpec) > 1 {
		return selection, deployOrder, l.Error(fmt.Errorf("%d ddo.cue file(s) found", len(actionSpec)))
	}

	l.Infof("Reading action specification %v", actionSpec[0])
	actionJson, e := container.WithExec(
		configuration.New(
			actionSpec[0],
			nil,
		).WithPackage("actions").AsJson(),
	).Stdout(ctx)

	actionJson = strings.TrimRight(actionJson, "\r\n")

	if e != nil {
		return selection, deployOrder, l.Error(e)
	}

	l.Infof("Get selection: %v", arg.ActionsPath())
	selection = gjson.Get(actionJson, "actions."+arg.ActionsPath()+"|@pretty")
	if !selection.Exists() {
		return selection, deployOrder, l.Error(fmt.Errorf("no such path: %v", arg.ActionsPath()))
	}

	if !gjson.Valid(selection.String()) {
		return selection, deployOrder, l.Error(fmt.Errorf("resulting json-selection from path is invalid"))
	}

	l.Infof("Get deployOrder")
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
		// in case of evomer and non #resourceId - not catching then signal, due to the if-break?
		// with a little more wait, the error signal is picked up.It's ok with va, if - slower stuff
		time.Sleep(25 * time.Millisecond)
		stopListening = true
		close(signalError)
	}

	if noOfErrors > 0 {
		return l.Error(fmt.Errorf("%v component(s) failed", noOfErrors))
	}
	return nil
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
