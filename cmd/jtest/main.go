package main

import (
	"ddo/alogger"
	"ddo/configuration"
	"fmt"
	"github.com/tidwall/gjson"
	"os"
	"strings"
)

var l = alogger.New() //.Caller()

//type Component struct {
//	Folder string   `json:"folder"`
//	Tags   []string `json:"tags"`
//}

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

	const (
		argProgramName = iota
		argPathToActionSpecification
		argActionsPath
		argMinNo = argActionsPath + 1
	)

	l.Infof("Start program: %v", os.Args[argProgramName])

	if len(os.Args) < argMinNo {
		return l.Error(
			fmt.Errorf(
				"missing parameter(s) - usage: PROGRAM <path to action specification> <actions path>"))
	}

	rawJson, err := configuration.New(os.Args[argPathToActionSpecification], nil).AsJson()

	if err != nil {
		return l.Error(fmt.Errorf("failure getting ddo actions: %v", err))
	}

	path := "actions." + strings.Join(os.Args[argActionsPath:], ".")
	l.Infof("Action path: %s", path)
	actions := gjson.GetBytes(rawJson, path+"|@pretty")
	if !actions.Exists() {
		return l.Error(fmt.Errorf("no such path: %v", path))
	}

	l.Infof("actions: \n%v ", actions)

	//deployOrder := gjson.GetBytes(rawJson, "deployOrder")

	// resolve path to slice of components
	//var components []Component
	_, ok := gjson.Parse(actions.String()).Value().(map[string]interface{})
	if !ok {
		return l.Error(fmt.Errorf("error parsing json"))
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

//func objectIsComponent(v map[string]interface{}) bool {
//	_, hasFolder := v["folder"]
//	_, hasTags := v["tags"]
//	return hasFolder && hasTags
//}
//
//func valueIsObject(v interface{}) (bool, map[string]interface{}) {
//	switch value := v.(type) {
//	case map[string]interface{}:
//		return true, value
//	default:
//		return false, nil
//	}
//}

//func printKeys(m map[string]interface{}, level, padding int) {
//	for k, v := range m {
//		if level == 1 {
//			fmt.Printf("%*s:\n", padding+len(k), k)
//		} else if level == 2 {
//			fmt.Printf("%*s: ", padding+len(k), k)
//		} else {
//			fmt.Printf("%s ", k)
//		}
//		isObject, value := valueIsObject(v)
//		if isObject && !objectIsComponent(value) {
//			printKeys(value, level+1, padding+2)
//		}
//	}
//	fmt.Println()
//}
