package deployment

import (
	"uuid"
)

// definition for a deployment
#deployment: {

	// bicep/arm file path, relative to repo root
	templatePath: string

	// template parameters
	parameters: [string]: _

	// target for the deployment
	target: #resourceGroup | #subscription | #managementGroup
}

// ResourceGroup target
#resourceGroup: {
	name:             string
	inSubscriptionId: #guid
}

// Subscription target
#subscription: {
	#base
}

// ManagementGroup target
#managementGroup: {
	#base
}

// common base for subscription and management group targets
#base: {
	id:       #guid
	location: string
}

// string as guid type
#guid: alias=uuid.Parse(alias)
