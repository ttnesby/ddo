package specification

import (
	"ddo/alogger"
	"ddo/arg"
)

var l alogger.ALogger

func Init() {
	l = alogger.New(arg.InDebugMode())
}

func GetComponentsInPath() (components []string) {
	return nil
}
