package ddospec

#operations: ["ce","va","if","de"]
#anOperation: or(#operations)

componentsPath: string
#component: {
	folder: string
	tags: [...string]
}
#components: [string]: #component
#componentActions: {...}
#listOfComponents: [...string]
componentOrder: [...#listOfComponents]

actions: {
	for o in #operations {"\(o)": _}
	[operation=#anOperation]: #componentActions
}