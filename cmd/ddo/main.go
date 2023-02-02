package main

import (
	"context"
	"ddo/action"
	"ddo/arg"
	"ddo/configuration"
	"ddo/deployment"
	"ddo/path"
	"os"
)

func init() {
	arg.Init()
	action.Init()
	configuration.Init()
	deployment.Init()
	path.Init()
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
