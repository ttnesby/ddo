package arg

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	flagArgActionsPath = 0
	flagArgMinNo       = flagArgActionsPath + 1

	usage = `usage: %s
ddo [options] [action path...]

Action path is a path to component in ddo.cue file, starting with one of:
ce - config export
va - validate config against azure
if - what-if analysis against azure
de - deploy to azure

e.g. 
ddo ce navutv rg - config export of navutv and component rg 
ddo ce navutv - config export of all components in navutv
ddo if - what-if of all components in all tenants

Options:
`
)

var config struct {
	debug    bool
	noResult bool
}

func Init() {
	flag.BoolVar(&config.debug, "debug", false, "debug mode")
	flag.BoolVar(&config.noResult, "no-result", false, "No display of action result")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
}

func AreOk() bool {
	switch {
	case len(flag.Args()) == 0:
		fmt.Printf("missing parameter(s)\n\n")
		flag.Usage()
		return false
	case flag.Arg(0) != "ce" && flag.Arg(0) != "va" && flag.Arg(0) != "if" && flag.Arg(0) != "de":
		fmt.Printf("invalid parameter(s)\n\n")
		flag.Usage()
		return false
	default:
		return true
	}
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
	return config.debug
}

func NoResultDisplay() bool {
	return config.noResult
}
