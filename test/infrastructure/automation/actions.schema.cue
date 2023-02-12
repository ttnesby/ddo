package actions

import (
	"strings"
)

//########## 3 mandatory fields ##########

// repo relative path to infrastructure components
_componentsPath: #nonEmptyString
// deploy order of infrastructure components
deployOrder: [...#listOfComponents]
// ce, va, if, de - and which additional configurations actions for infrastructure components should be possible
actions: #actions

//########## schema ##########

#operations: ["ce", "va", "if", "de", "evomer"]
#anOperation: or(#operations)

#nonEmptyString: s={
	string
	_len:   len(s)
	#valid: >0 & _len
}

#component: {
	folder: #folderPath
	tags: [...string]
}

#folderPath: s={
	string
	#valid: strings.HasPrefix(s, _componentsPath+"/") & true
}

#components: close({[string]: #component})
#componentNames: [ for k, _ in #components {k}]
#aComponentName: or(#componentNames)

#listOfComponents: [#aComponentName, ...#aComponentName]

#actions: {
	#componentActions: {...}

	for o in #operations {"\(o)": _}
	[operation=#anOperation]: #componentActions
	... // extendable
}
