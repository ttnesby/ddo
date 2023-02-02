package arg

import (
	"ddo/alogger"
	"os"
	"strings"
)

var l = alogger.New()

const (
	argProgramName = iota
	argActionsPath
	argMinNo = argActionsPath + 1
)

func AreOk() bool {
	l.Infof("Start program: %v", os.Args[argProgramName])

	if len(os.Args) < argMinNo {
		l.Errorf("missing parameter(s)")
		l.Infof(`\n
usage: ddo <operation> <actions path...>

<operation> - one of: ce, va, if, de
ce - config export
va - validate config against azure
if - what-if analysis against azure
de - deploy to azure

<actions path...> - path to component in ddo.cue file, 

e.g. 
- "ddo ce navutv rg" for config export of navutv and component rg 
- "ddo ce navutv" for config export of all components in navutv
- "ddo if" for what-if of all components in all tenants
`)
		return false
	}

	return true
}

func Operation() string {
	return os.Args[argActionsPath]
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
