package arg

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	argProgramName     = 0
	flagArgActionsPath = 0
	flagArgMinNo       = flagArgActionsPath + 1
)

var (
	debug    bool
	noResult bool
)

func Init() {
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.BoolVar(&noResult, "noResult", false, "No display of action result")
	flag.Parse()
}

func AreOk() bool {

	fmt.Printf("Start program: %v\n", os.Args[argProgramName])

	if len(os.Args) < flagArgMinNo {
		fmt.Printf("missing parameter(s)")
		fmt.Printf(`\n
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
\n`)
		return false
	}

	return true
}

func Operation() string {
	return flag.Args()[flagArgActionsPath]
}

func ActionsPath() (actionPath string) {
	return strings.Join(flag.Args()[flagArgActionsPath:], ".")
}

func LastActions() []string {
	if len(flag.Args()) == flagArgMinNo {
		return []string{}
	} else {
		return flag.Args()[flagArgMinNo:]
	}
}

func InDebugMode() bool {
	return debug
}

func NoResultDisplay() bool {
	return noResult
}
