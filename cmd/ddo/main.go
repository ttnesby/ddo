package main

import (
	"context"
	"ddo/action"
	"ddo/action/component"
	"ddo/arg"
	delete2 "ddo/azcli/delete"
	"ddo/azcli/deployment"
	"ddo/configuration"
	"ddo/path"
	"os"
)

func init() {
	arg.Init()
	action.Init()
	component.Init()
	configuration.Init()
	deployment.Init()
	path.Init()
	delete2.Init()
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
