package resource

import (
	"ddo/alogger"
	"ddo/arg"
)

var l alogger.ALogger

func Init() {
	l = alogger.New(arg.InDebugMode())
}

type AzCli []string

func Delete(rId string) (azCmd AzCli) {

	azCmd = []string{
		"az",
		"resource",
		"delete",
		"--ids",
		rId,
		"--verbose",
	}

	l.Debugf("azCmd: %v", azCmd)
	return azCmd
}

func Show(rId string) (azCmd AzCli) {

	azCmd = []string{
		"az",
		"resource",
		"show",
		"--ids",
		rId,
	}

	l.Debugf("azCmd: %v", azCmd)
	return azCmd
}
