package actions

// repo relative path to infrastructure components
componentsPath: #componentsPath
// deploy order of infrastructure components
deployOrder: [...#listOfComponents]
// ce, va, if, de - and which additional configurations actions for infrastructure components should be possible
actions: [string]: _ //use #actions

#operations: ["ce","va","if","de"]
#anOperation: or(#operations)

#componentsPath: s={
	string
	_len: len(s)
	#valid: > 0 & _len
}

#component: {
	folder: string
	tags: [...string]
}

#components: close({[string]: #component})
#componentNames: [for k,_ in #components {k}]
#aComponentName: or(#componentNames)

#listOfComponents: [#aComponentName,...#aComponentName]

#actions: {
	#componentActions: {...}
	actions: {
			for o in #operations {"\(o)": _}
			[operation=#anOperation]: #componentActions
	}
}