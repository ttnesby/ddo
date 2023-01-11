package level

type Level string

const (
	ResourceGroup   Level = "group"
	Subscription    Level = "sub"
	ManagementGroup Level = "mg"
)

func Levels() []Level {
	return []Level{ResourceGroup, Subscription, ManagementGroup}
}

func (l Level) Valid() bool {
	switch l {
	case ResourceGroup, Subscription, ManagementGroup:
		return true
	default:
		return false
	}
}
