package deployment

import (
	"encoding/json"
	"encoding/base64"
	g "ddo.test/test:global"
)

// cue export ./test/infrastructure/resourceGroup -t tenant=navutv

_tenant:          g.#aTenantKey @tag(tenant)
_crid:             *"" | string  @tag(crid)
_crpoltruststatus: *"" | string  @tag(crpoltruststatus)
_crb64:            *"" | string  @tag(crb64)

// load the complete container registry data (az resource show) or empty object
_crData: [
	if len(_crb64) > 0 {json.Unmarshal("\(base64.Decode(null, _crb64))")},
	{}
	][0]

// using the existing container registry data from azure, for e.g. tag data injection
_tagCRId: [
	if len(_crid) > 0 { "CRId": _crid },
	{}
	][0]
_tagCRTrustStatus: [
	if len(_crpoltruststatus) > 0 { "CRTrustStatus": _crpoltruststatus },
	{}
	][0]
_tagAnonymousPullEnabled: [
	if len(_crData) > 0 { "AnonymousPullEnabled": _crData.properties.anonymousPullEnabled },
	{}
	][0]

#name:     "container-registry"
#location: g.#location.norwayeast
#tags:     {...,g.#tagsTemplate} & _tagCRId & _tagCRTrustStatus & _tagAnonymousPullEnabled

// if you want data injection from this component or "evomer" it with ddo, define the resource id
#resourceId: "/subscriptions/\(g.#subscriptionId[_tenant])/resourceGroups/\(#name)"

templatePath: "./test/infrastructure/resourceGroup/main.bicep"

parameters: #jsonParameterFile & {
	#s: {
		name:     #name
		location: #location
		tags:     #tags
	}
}

target: subscription: {
	id:       g.#subscriptionId[_tenant]
	location: g.#location.norwayeast
}
