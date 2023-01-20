package actions

// cue export ./test/infrastructure/automation

#tenants: ["navutv", "navno"]
#aTenant: or(#tenants)

componentsPath: "./test/infrastructure"
//componentsPath: ""

#components: {
	#tenant: #aTenant
	rg: #component & {
		folder: "resourceGroup"
		tags: ["tenant=\(#tenant)"]
	}
	cr: #component &{
		folder: "containerRegistry"
		tags:["tenant=\(#tenant)"]
	}
}

deployOrder: [["rg"], ["cr"]]

_actions: {
	for t in #tenants {"\(t)": _}
	[tenant=#aTenant]: #components & {#tenant: tenant}
}

actions: #actions & {#componentActions: _actions}