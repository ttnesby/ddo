package actions

// cue export ./test/infrastructure/automation

import (
	g "ddo.test/test/infrastructure:global"
)

_componentsPath: "./test/infrastructure"

#components: {
	#tenant: g.#aTenantKey
	rg: #component & {
		folder: "\(_componentsPath)/resourceGroup"
		tags: ["tenant=\(#tenant)"]
	}
	cr: #component &{
		folder: "\(_componentsPath)/containerRegistry"
		tags:["tenant=\(#tenant)"]
	}
}

deployOrder: [["rg"], ["cr"]]

_actions: {
	for t in g.#tenants {"\(t)": _}
	[tenant=g.#aTenantKey]: #components & {#tenant: tenant}
}

actions: #actions & {#componentActions: _actions}