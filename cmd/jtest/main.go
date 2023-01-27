package main

import (
	"ddo/configuration"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
	"os"
	"strings"
	"time"
)

type Component struct {
	Folder string   `json:"folder"`
	Tags   []string `json:"tags"`
}

func main() {
	exitCode := func() int {
		if err := do(); err != nil {
			return 1
		}
		return 0
	}()
	os.Exit(exitCode)
}

func do() error {

	logError := func(e error) error {
		log.Err(e).Msg("Error")
		return e
	}

	const (
		argProgramName = iota
		argPathToActionSpecification
		argActionsPath
	)

	log.Logger = (log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})).With().Caller().Logger()

	log.Info().Str("program", os.Args[argProgramName]).Msg("Start")
	if len(os.Args) < argActionsPath+1 {
		return logError(
			fmt.Errorf(
				"missing parameter(s) - usage: PROGRAM <path to action specification> <actions path>"))
	}

	rawJson, err := configuration.New(os.Args[argPathToActionSpecification], nil).AsJson()

	if err != nil {
		return logError(fmt.Errorf("failure getting ddo actions: %v", err))
	}

	path := "actions." + strings.Join(os.Args[argActionsPath:], ".")
	log.Info().Str("Action path", path).Msg("Resolve")
	actions := gjson.GetBytes(rawJson, path+"|@pretty")
	if !actions.Exists() {
		return logError(fmt.Errorf("no such path: %v", path))
	}

	fmt.Printf("actions: %v \n", actions)

	//deployOrder := gjson.GetBytes(rawJson, "deployOrder")

	// resolve path to slice of components
	//var components []Component
	_, ok := gjson.Parse(actions.String()).Value().(map[string]interface{})
	if !ok {
		return logError(fmt.Errorf("error parsing json"))
	}

	//for k, v := range m {
	//	fmt.Printf("key: %v, value: %v\n", k, v)
	//}

	//var plan Plan
	//err = json.Unmarshal(rawJson, &plan)
	//if err != nil {
	//	_ = fmt.Errorf("error unmarshalling json: %v", err)
	//	os.Exit(1)
	//}
	//
	//fmt.Printf("componentOrder: %v \n", plan.DeployOrder)
	//
	//fmt.Printf("Available actions: \n")
	//printKeys(plan.Actions, 1, 2)
	//
	return nil
}

func objectIsComponent(v map[string]interface{}) bool {
	_, hasFolder := v["folder"]
	_, hasTags := v["tags"]
	return hasFolder && hasTags
}

func valueIsObject(v interface{}) (bool, map[string]interface{}) {
	switch value := v.(type) {
	case map[string]interface{}:
		return true, value
	default:
		return false, nil
	}
}

func printKeys(m map[string]interface{}, level, padding int) {
	for k, v := range m {
		if level == 1 {
			fmt.Printf("%*s:\n", padding+len(k), k)
		} else if level == 2 {
			fmt.Printf("%*s: ", padding+len(k), k)
		} else {
			fmt.Printf("%s ", k)
		}
		isObject, value := valueIsObject(v)
		if isObject && !objectIsComponent(value) {
			printKeys(value, level+1, padding+2)
		}
	}
	fmt.Println()
}
