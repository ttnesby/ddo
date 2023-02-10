package main

import (
	"context"
	"ddo/action"
	"ddo/action/component"
	"ddo/arg"
	"ddo/azcli/deployment"
	"ddo/azcli/resource"
	"ddo/cuecli"
	"ddo/path"
	"os"
)

func init() {
	arg.Init()
	action.Init()
	component.Init()
	cuecli.Init()
	deployment.Init()
	path.Init()
	resource.Init()
}

func main() {

	exitCode := func() int {
		switch arg.AreOk() {
		case true:
			if err := action.Do(context.Background()); err != nil {
				return 1
			}
			return 0
		default:
			return 1
		}
	}()
	os.Exit(exitCode)
}
