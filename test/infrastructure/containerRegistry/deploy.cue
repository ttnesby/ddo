package deployment

import (
	g "ddo.test/test:global"
	rg "ddo.test/test/infrastructure/resourceGroup:deployment"
	ddo "github.com/ttnesby/ddoapi/cue/v1:deployment"
)

// cue export -p deployment ./test/infrastructure/containerRegistry ./cue.mod/pkg/github.com/ttnesby/ddoapi/cue/v1/deployment.schema.cue -t tenant=navutv

_tenant: g.#aTenantKey @tag(tenant)

#name:     "azplatform\(_tenant)"
#location: g.#location.norwayeast
#tags:     g.#tagsTemplate
#skuName:  "Basic" | "Classic" | "Standard" | *"Premium"

#enabling: "enabled" | "disabled"

#propertiesTemplate: {
	adminUserEnabled:     *false | bool
	anonymousPullEnabled: *true | bool // make it easy to pull images with no security concerns
	dataEndpointEnabled:  *false | bool
	encryption: status: *"disabled" | #enabling
	networkRuleBypassOptions: *"None" | "AzureServices"
	policies: {
		azureADAuthenticationAsArmPolicy: *"enabled" | #enabling
		exportPolicy:                     *"enabled" | #enabling
		quarantinePolicy:                 *"disabled" | #enabling
		retentionPolicy: {// delete older images not used
			days:   *5 | int
			status: *"enabled" | #enabling
		}
		softDeletePolicy: status: *"disabled" | #enabling
		trustPolicy: status:      *"disabled" | #enabling
	}
	publicNetworkAccess: *"enabled" | #enabling // must be enabled for sku Standard - not catched by validate?!
	zoneRedundancy:      *"enabled" | #enabling // requires sku Premium
}

templatePath: "./test/containerRegistry/main.bicep"

parameters: ddo.#jsonParameterFile & {
		#s:  {
			name:     #name
			location: #location
			tags:     #tags
			skuName:  #skuName
			// no identity
			properties: #propertiesTemplate
		}
}

target: resourceGroup: {
		name:             rg.#name
		inSubscriptionId: g.#subscriptionId[_tenant]
}

