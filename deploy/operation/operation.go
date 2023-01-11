package operation

type Operation string

const (
	Validate Operation = "validate"
	WhatIf   Operation = "create --what-if"
	Deploy   Operation = "create"
)

func Operations() []Operation {
	return []Operation{Validate, WhatIf, Deploy}
}
