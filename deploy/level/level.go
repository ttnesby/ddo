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
