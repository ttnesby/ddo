package deployment

import (
	g "ddo.test/test:global"
	ddo "github.com/ttnesby/ddoapi/cue/v1:deployment"
)

// cue export -p deployment ./test/infrastructure/resourceGroup ./cue.mod/pkg/github.com/ttnesby/ddoapi/cue/v1/deployment.schema.cue -t tenant=navutv

_tenant: g.#aTenantKey @tag(tenant)

#name:     "container-registry"
#location: g.#location.norwayeast
#tags:     g.#tagsTemplate

templatePath: "./infrastructure/resourceGroup/main.bicep"

parameters: ddo.#jsonParameterFile & {
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
