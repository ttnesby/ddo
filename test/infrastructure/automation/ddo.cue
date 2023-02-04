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
		// TODO: add support for data dependencies
		// support of data dependencies iff another resource exists in azure
		// due to cyclic dependencies and order of creation, need deploy time injection using tags
		// e.g. this component depends cr component, which depends on this component...

		// two scenarios:
		// 1. require a specific field from cr (az resource show --ids <cr id>) + gjson lookup
		// tags:[..., "crid=cr.id","crpoltruststatus=cr.policies.trustPolicy.status"]
		// string only - either data or ""

		// 2. require complex data (list, struct) from cr, base64 encode
		// tags:[..., crb=cr.b64"]
		// [] or {} - either data or {},[]
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