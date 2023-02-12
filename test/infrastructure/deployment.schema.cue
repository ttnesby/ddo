package deployment

import (
	"uuid"
)

//########## 3 mandatory fields ##########

// bicep/arm file path, relative to repo root
templatePath: #nonEmptyString
// template parameters
parameters: #jsonParameterFile
// target for the deployment
target: #resourceGroup | #subscription | #managementGroup

//########## schema ##########

#nonEmptyString: s={
	string
	_len:   len(s)
	#valid: >0 & _len
}

// ResourceGroup target
#resourceGroup: resourceGroup: {
	name:             #nonEmptyString
	inSubscriptionId: #guid
}

// Subscription target
#subscription: subscription: #base

// ManagementGroup target
#managementGroup: managementGroup: #base

// common properties for subscription and management group targets
#base: {
	id:       #guid
	location: #nonEmptyString
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
