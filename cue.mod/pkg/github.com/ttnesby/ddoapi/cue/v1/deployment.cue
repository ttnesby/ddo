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

#jsonParameterFile: {
	_s=#s: {...}
	"$schema":        "https://schema.management.azure.com/schemas/2015-01-01/deploymentParameters.json#"
	"contentVersion": "1.0.0.0"
	parameters:       #value & {#s: _s}
}

#value: {
	#s: {...}

	if len(#s) == 0 {{}}
	if len(#s) > 0 {for k, v in #s {
		if (v.reference != _|_) {"\(k)": v}
		if (v.reference == _|_) {"\(k)": value: v}
	}}
}