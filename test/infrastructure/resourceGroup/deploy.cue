package deployment

import (
	g "ddo.test/test:global"
)

// cue export ./test/infrastructure/resourceGroup -t tenant=navutv

_tenant: g.#aTenantKey @tag(tenant)

#name:     "container-registry"
#location: g.#location.norwayeast
#tags:     g.#tagsTemplate

templatePath: "./infrastructure/resourceGroup/main.bicep"

parameters: #jsonParameterFile & {
		#s: {
				name:     #name
				location: #location
				tags:     #tags
		}
}

target: subscription : {
		id:       g.#subscriptionId[_tenant]
		location: g.#location.norwayeast
}
