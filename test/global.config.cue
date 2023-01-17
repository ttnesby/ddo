package global

import (
	"uuid"
)

#guid: alias=uuid.Parse(alias)

#tenantId: {
	navno:  "62366534-1ec3-4962-8869-9b5535279d0b"
	navutv: "82bdf6c1-3e56-4a5e-8c50-c331165e0192"
}
#tenantId: close({[string]: #guid})
#tenants: [ for k, _ in #tenantId {k}]
#aTenantKey: or(#tenants)

// ./nav/platform/management
#subscriptionId: {
	navno:  "7e260459-3026-4653-b259-0347c0bb5970"
	navutv: "ca1e4592-6211-4c03-aac5-e681c1d1ea0d"
}
#subscriptionId:  close({[#aTenantKey]: #guid})
#aSubscriptionId: or([ for _, v in #subscriptionId {v}])

// gha_platform
#appRegistrationClientId: {
	navno:  "6afc2353-8385-414b-b05b-730f467ce37a"
	navutv: "a82864fa-5cef-453f-989a-a722242f81ba"
}
#appRegistrationClientId: close({[#aTenantKey]: #guid})

#location: {
	norwayeast: "norwayeast"
}
#location:       close({[string]: string})
#aLocationValue: or([ for _, v in #location {v}])

#tagsTemplate: {
	Team:              *"azure-platform" | string
	"Technical Owner": *"azure-platform" | string
	Provenance:        "https://github.com/ttnesby/ddo"
	Application:       "ddo"
}
