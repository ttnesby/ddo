package main

import (
	"encoding/json"
	"fmt"
	"os"
)

//type Component struct {
//	Folder string   `json:"folder"`
//	Tags   []string `json:"tags"`
//}

func main() {
	rawJson := `
{
    "componentsPath": "./test/infrastructure",
    "componentOrder": [
        [
            "rg"
        ],
        [
            "cr"
        ]
    ],
    "actions": {
        "ce": {
            "navutv": {
                "rg": {
                    "folder": "resourceGroup",
                    "tags": [
                        "tenant=navutv"
                    ]
                },
                "cr": {
                    "folder": "containerRegistry",
                    "tags": [
                        "tenant=navutv"
                    ]
                }
            },
            "navno": {
                "rg": {
                    "folder": "resourceGroup",
                    "tags": [
                        "tenant=navno"
                    ]
                },
                "cr": {
                    "folder": "containerRegistry",
                    "tags": [
                        "tenant=navno"
                    ]
                }
            }
        },
        "va": {
            "navutv": {
                "rg": {
                    "folder": "resourceGroup",
                    "tags": [
                        "tenant=navutv"
                    ]
                },
                "cr": {
                    "folder": "containerRegistry",
                    "tags": [
                        "tenant=navutv"
                    ]
                }
            },
            "navno": {
                "rg": {
                    "folder": "resourceGroup",
                    "tags": [
                        "tenant=navno"
                    ]
                },
                "cr": {
                    "folder": "containerRegistry",
                    "tags": [
                        "tenant=navno"
                    ]
                }
            }
        },
        "if": {
            "navutv": {
                "rg": {
                    "folder": "resourceGroup",
                    "tags": [
                        "tenant=navutv"
                    ]
                },
                "cr": {
                    "folder": "containerRegistry",
                    "tags": [
                        "tenant=navutv"
                    ]
                }
            },
            "navno": {
                "rg": {
                    "folder": "resourceGroup",
                    "tags": [
                        "tenant=navno"
                    ]
                },
                "cr": {
                    "folder": "containerRegistry",
                    "tags": [
                        "tenant=navno"
                    ]
                }
            }
        },
        "de": {
            "navutv": {
                "rg": {
                    "folder": "resourceGroup",
                    "tags": [
                        "tenant=navutv"
                    ]
                },
                "cr": {
                    "folder": "containerRegistry",
                    "tags": [
                        "tenant=navutv"
                    ]
                }
            },
            "navno": {
                "rg": {
                    "folder": "resourceGroup",
                    "tags": [
                        "tenant=navno"
                    ]
                },
                "cr": {
                    "folder": "containerRegistry",
                    "tags": [
                        "tenant=navno"
                    ]
                }
            }
        }
    }
}
`

	//type Operations struct {
	//	Ce map[string]interface{} `json:"ce"`
	//	Va map[string]interface{} `json:"va"`
	//	If map[string]interface{} `json:"if"`
	//	De map[string]interface{} `json:"de"`
	//}
	type Plan struct {
		ComponentsPath string     `json:"componentsPath"`
		ComponentOrder [][]string `json:"componentOrder"`
		//Actions        Operations `json:"actions"`
		Actions map[string]interface{} `json:"actions"`
	}

	var plan Plan
	err := json.Unmarshal([]byte(rawJson), &plan)
	if err != nil {
		_ = fmt.Errorf("error unmarshalling json: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Arguments: %v \n\n", os.Args[1:])

	fmt.Printf("componentsPath: %s \n", plan.ComponentsPath)
	fmt.Printf("componentOrder: %v \n", plan.ComponentOrder)

	fmt.Printf("Available actions: \n")
	printKeys(plan.Actions, 1, 2)

	os.Exit(0)
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
			fmt.Printf("%*s\n", padding+len(k), k)
		} else if level == 2 {
			fmt.Printf("%*s - ", padding+len(k), k)
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
