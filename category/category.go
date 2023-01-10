package category

import "fmt"

type Category int

const (
	Validate Category = iota
	WhatIf
	Deploy
)

func Categories() []Category {
	return []Category{Validate, WhatIf, Deploy}
}

func (c Category) ToOperation() (string, error) {

	catToOp := map[Category]string{
		Validate: "validate",
		WhatIf:   "create --what-if",
		Deploy:   "create",
	}

	op, ok := catToOp[c]

	if !ok {
		return "", fmt.Errorf("category %v not found", c)
	}
	return op, nil
}
