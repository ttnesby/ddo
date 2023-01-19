package ddospec

#tenants: ["navutv", "navno"]
#aTenant: or(#tenants)

componentsPath: "./test/infrastructure"
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

#componentActions: {
	for t in #tenants {"\(t)": _}
	[tenant=#aTenant]: #components & {#tenant: tenant}
}

componentOrder: [["rg"], ["cr"]]
