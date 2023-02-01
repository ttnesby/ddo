package main

import (
	"context"
	"ddo/action"
	"ddo/arg"
	"os"
)

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
