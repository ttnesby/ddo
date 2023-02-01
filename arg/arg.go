package arg

import (
	"ddo/alogger"
	"ddo/path"
	"os"
	"strings"
)

var l = alogger.New()

const (
	argProgramName = iota
	argPathToActionSpecification
	argActionsPath
	argMinNo = argActionsPath + 1
)

func AreOk() bool {
	l.Infof("Start program: %v", os.Args[argProgramName])
	l.Debugf("Check path %v", os.Args[argPathToActionSpecification])

	if !path.AbsExists(path.RepoAbs(os.Args[argPathToActionSpecification])) {
		l.Errorf("cannot find %v", os.Args[argPathToActionSpecification])
		return false
	}

	if len(os.Args) < argMinNo {
		l.Errorf("missing parameter(s) - usage: PROGRAM <path to action specification> <actions path...>")
		return false
	}

	return true
}

func ActionSpecification() (relativePath string) {
	return os.Args[argPathToActionSpecification]
}

func ActionsPath() (actionPath string) {
	return strings.Join(os.Args[argActionsPath:], ".")
}

func LastActions() []string {
	if len(os.Args) == argMinNo {
		return []string{}
	} else {
		return os.Args[argMinNo:]
	}
}
