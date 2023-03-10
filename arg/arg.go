package arg

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	OpCE = "ce"
	OpVA = "va"
	OpIF = "if"
	OpDE = "de"
	OpRE = "evomer"

	flagArgActionsPath = 0
	flagArgMinNo       = flagArgActionsPath + 1

	usage = `usage: %s
ddo [options] [action path...]

Action path is a path to component in ddo.cue file, starting with one of:
ce - config export
va - validate config against azure
if - what-if analysis against azure
de - deploy to azure
evomer - remove iff the component has '#resourceId' definition

e.g. 
ddo ce navutv rg - config export of navutv and component rg 
ddo ce navutv - config export of all components in navutv
ddo if - what-if of all components in all tenants

Options:
`
)

var config struct {
	debug          bool
	debugContainer bool
	noResult       bool
	containerRef   string
}

func Init() {
	flag.BoolVar(&config.debug, "debug", false, "debug mode")
	flag.BoolVar(&config.debugContainer, "debug-container", false, "debug mode for dagger.io")
	flag.BoolVar(&config.noResult, "no-result", false, "No display of action result")
	flag.StringVar(&config.containerRef, "cnt", "docker.io/ttnesby/azbicue:latest",
		"container ref. hosting az cli, bicep and cue. Default is 'docker.io/ttnesby/azbicue:latest'")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
}

func AreOk() bool {
	switch {
	case len(flag.Args()) == 0:
		fmt.Printf("missing action path\n\n")
		flag.Usage()
		return false
	case flag.Arg(0) != OpCE &&
		flag.Arg(0) != OpVA &&
		flag.Arg(0) != OpIF &&
		flag.Arg(0) != OpDE &&
		flag.Arg(0) != OpRE:
		fmt.Printf("invalid action path\n\n")
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

func DebugContainer() bool {
	return config.debugContainer
}

func NoResultDisplay() bool {
	return config.noResult
}

func ContainerRef() string {
	return config.containerRef
}
