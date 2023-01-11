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

func (o Operation) Valid() bool {
	switch o {
	case Validate, WhatIf, Deploy:
		return true
	default:
		return false
	}
}
