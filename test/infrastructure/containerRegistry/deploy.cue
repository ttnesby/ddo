package containerRegistry

import (
	g "ddo.test:global"
	rg "ddo.test/infrastructure/resourceGroup"
	ddo "github.com/ttnesby/ddoapi/cue/v1:deployment"
)

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

ddo.#deployment & {

	templatePath: "./test/containerRegistry/main.bicep"
	parameters: {
		name:     #name
		location: #location
		tags:     #tags
		skuName:  #skuName
		// no identity
		properties: #propertiesTemplate
	}

	target: ddo.#resourceGroup & {
		name: rg.#name
		inSubscriptionId:  g.#subscriptionId[_tenant]
	}
}
