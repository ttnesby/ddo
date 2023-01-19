package resourceGroup

import (
	g "ddo.test/test:global"
	ddo "github.com/ttnesby/ddoapi/cue/v1:deployment"
)

_tenant: g.#aTenantKey @tag(tenant)

#name:     "container-registry"
#location: g.#location.norwayeast
#tags:     g.#tagsTemplate

ddo.#deployment & {

	templatePath: "./infrastructure/resourceGroup/main.bicep"
	parameters: ddo.#jsonParameterFile & {
		#s:  {
				name:     #name
				location: #location
				tags:     #tags
		}
	}

	target: ddo.#subscription & {
		id:       g.#subscriptionId[_tenant]
		location: g.#location.norwayeast
	}
}
