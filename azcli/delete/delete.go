package delete

import (
	"ddo/alogger"
	"ddo/arg"
)

var l alogger.ALogger

func Init() {
	l = alogger.New(arg.InDebugMode())
}

type AzCli []string

func ResourceId(rId string) (azCmd AzCli) {

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
