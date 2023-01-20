package actions

import (
	ddo "github.com/ttnesby/ddoapi/cue/v1:actions"
)

// cue export -p actions ./test/ddo.cue ./cue/actions.schema.cue

#tenants: ["navutv", "navno"]
#aTenant: or(#tenants)

componentsPath: "./test/infrastructure"
//componentsPath: ""

#components: {
	#tenant: #aTenant
	rg: ddo.#component & {
		folder: "resourceGroup"
		tags: ["tenant=\(#tenant)"]
	}
	cr: ddo.#component &{
		folder: "containerRegistry"
		tags:["tenant=\(#tenant)"]
	}
}

_actions: {
	for t in #tenants {"\(t)": _}
	[tenant=#aTenant]: #components & {#tenant: tenant}
}

ddo.#actions & {#componentActions: _actions}

deployOrder: [["rg"], ["cr"]]