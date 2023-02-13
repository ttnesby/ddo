# Infrastructure as Configuration

`ddo` (DaggerDO) is a `thin` concept tool for `azure infrastructure as configuration`. The building blocks are:
> - Use [cue language](https://cuelang.org/) for json parameters and deployment destination for bicep templates
> - Use cue language for orchestrating deployment of multiple templates, including order and data dependencies
> - Use of [dagger.io](https://dagger.io/) for supporting both local - and GitHub actions CI/CD
> - Use of an OCI image with three basic tools: `az cli`, `bicep cli` and `cue cli`

The core actions for one or multiple infrastructure components are:
> - `ce` - export configuration (`cue export`)
> - `va` - validate configuration (`az deployment <dest> validate`)
> - `if` - what-if of configuration (`az deployment <dest> create --what-if`)
> - `de` - deploy configuration (`az deployment <dest> create`)
> - `evomer` - remove configuration (`az resource remove --ids <ids>`)

`va` will report any regional restrictions or policy violations.

`if` is a dry-run of `de` and will report any changes to the infrastructure. Since each bicep template is simple in 
terms of just one template, the what-if is is working as good-as-it-can. The `false-positives`, depending on the type of 
infrastructure component, are still popping up.

> `ddo` is `path oriented`, meaning that the chosen action will be applied to all components in the path. See 
> [Orchestration configuration](#orchestration-configuration) for more details.

## Configuration example

`./test/infrastructure` is an example of a simple infrastructure solution. A resource group with a container registry, 
to be deployed to two different tenants. The `navutv` tenant is a `dev` tenant, while the `navno` tenant is a `prod` tenant.

>**Warning** cuelang versus json/yaml is a involving topic beyond the scope of this document. Briefly, it's a 
> configuration language with schema and import support included. Thus, configure once and 
> reuse where appropriate.

### Global configuration
- `./test/global.config.cue` is the global configuration for the solution, containing basic azure details
  - tenant guids
  - subscription guid
  - client id guid
  - location
  - relevant tags

The global configuration is used by the `resourceGroup` and `containerRegistry` components. In addition, it can be reused
in `Orchestration configuration` and in Github actions, e.g. federated login.

### Resource group configuration
- `./test/infrastructure/resourceGroup` is the resource group component
  - `./test/infrastructure/resourceGroup/main.bicep` is the bicep template
  - `./test/infrastructure/resourceGroup/deploy.cue` is the deployment configuration, defining just three elements
    - the relevant bicep template
    - parameters for the bicep template
    - the deployment destination
  
A more familiar look of the resourceGroup data for tenant `navutv`
```zsh
docker run -it --rm --platform linux/arm64 -v $(pwd):/rr --workdir="/rr" docker.io/ttnesby/tooling:latest
029f926b34d1:/rr# cue export ./test/infrastructure/resourceGroup -t tenant=navutv
# output
{
    "templatePath": "./test/infrastructure/resourceGroup/main.bicep",
    "parameters": {
        "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentParameters.json#",
        "contentVersion": "1.0.0.0",
        "parameters": {
            "name": {
                "value": "container-registry"
            },
            "location": {
                "value": "norwayeast"
            },
            "tags": {
                "value": {
                    "Team": "azure-platform",
                    "Technical Owner": "azure-platform",
                    "Provenance": "https://github.com/ttnesby/ddo",
                    "Application": "ddo"
                }
            }
        }
    },
    "target": {
        "subscription": {
            "id": "ca1e4592-6211-4c03-aac5-e681c1d1ea0d",
            "location": "norwayeast"
        }
    }
}
```
By changing `tenant=navutv` to `tenant=navno`, the output will be for that tenant. Configure once and reuse where 
appropriate makes life easier when working with more complex infrastructure solutions. No more copy-paste between 
multiple json parameter files.

### Container registry configuration
- `./test/infrastructure/containerRegistry` is the container registry component
  - `./test/infrastructure/containerRegistry/main.bicep` is the bicep template
  - `./test/infrastructure/resourceGroup/deploy.cue` is the deployment configuration (template, parameters and destination)

The `json-look`
```zsh
029f926b34d1:/rr# cue export ./test/infrastructure/containerRegistry -t tenant=navutv
# output
{
    "templatePath": "./test/infrastructure/containerRegistry/main.bicep",
    "parameters": {
        "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentParameters.json#",
        "contentVersion": "1.0.0.0",
        "parameters": {
            "name": {
                "value": "azplatformnavutv"
            },
            "location": {
                "value": "norwayeast"
            },
            "tags": {
                "value": {
                    "Team": "azure-platform",
                    "Technical Owner": "azure-platform",
                    "Provenance": "https://github.com/ttnesby/ddo",
                    "Application": "ddo"
                }
            },
            "skuName": {
                "value": "Premium"
            },
            "properties": {
                "value": {
                    "adminUserEnabled": false,
                    "anonymousPullEnabled": true,
                    "dataEndpointEnabled": false,
                    "encryption": {
                        "status": "disabled"
                    },
                    "networkRuleBypassOptions": "None",
                    "policies": {
                        "azureADAuthenticationAsArmPolicy": "enabled",
                        "exportPolicy": "enabled",
                        "quarantinePolicy": "disabled",
                        "retentionPolicy": {
                            "days": 5,
                            "status": "enabled"
                        },
                        "softDeletePolicy": {
                            "status": "disabled"
                        },
                        "trustPolicy": {
                            "status": "disabled"
                        }
                    },
                    "publicNetworkAccess": "enabled",
                    "zoneRedundancy": "enabled"
                }
            }
        }
    },
    "target": {
        "resourceGroup": {
            "name": "container-registry",
            "inSubscriptionId": "ca1e4592-6211-4c03-aac5-e681c1d1ea0d"
        }
    }
}
```

### Orchestration configuration
- `./test/infrastructure/automation`
  - `./test/infrastructure/automation/actions.schema.cue` schema and baseline for actions (ce, va, etc.)
  - `./test/infrastructure/automation/ddo.cue` configures the orchestration, defining just two elements
    - actions - relevant groupings under actions and components 
    - deployOrder - list of lists, defining the order of component deployment

The `json-look` with action `ce` only. The json for the other actions are similar.
```zsh
029f926b34d1:/rr# cue export ./test/infrastructure/automation
#output
{
    "deployOrder": [
        [
            "rg"
        ],
        [
            "cr"
        ]
    ],
    "actions": {
        "ce": {
            "navno": {
                "rg": {
                    "folder": "./test/infrastructure/resourceGroup",
                    "tags": [
                        "tenant=navno",
                        "crid=\u003c\u003cnavno.cr\u003c\u003cid",
                        "crpoltruststatus=\u003c\u003cnavno.cr\u003c\u003cproperties.policies.trustPolicy.status",
                        "crb64=\u003c\u003cnavno.cr\u003c\u003cb64"
                    ]
                },
                "cr": {
                    "folder": "./test/infrastructure/containerRegistry",
                    "tags": [
                        "tenant=navno"
                    ]
                }
            },
            "navutv": {
                "rg": {
                    "folder": "./test/infrastructure/resourceGroup",
                    "tags": [
                        "tenant=navutv",
                        "crid=\u003c\u003cnavutv.cr\u003c\u003cid",
                        "crpoltruststatus=\u003c\u003cnavutv.cr\u003c\u003cproperties.policies.trustPolicy.status",
                        "crb64=\u003c\u003cnavutv.cr\u003c\u003cb64"
                    ]
                },
                "cr": {
                    "folder": "./test/infrastructure/containerRegistry",
                    "tags": [
                        "tenant=navutv"
                    ]
                }
            }
        }
  }
}
```

By looking at the above json-file, the meaning of `path oriented` is easy to understand, it's a kind of path in the json file:
- `ddo ce navutv rg` will export the config for `resourceGroup` with `tenant=navutv` tag
- `ddo ce navutv cr` will export the config for `containerRegistry` with `tenant=navutv` tag
- `ddo ce navno` will export the config for both `resourceGroup` and `containerRegistry` with `tenant=navno` tag
- `ddo ce` will export the config for both components and both tenants

The same behavior is also available for the other actions. The relevant structure and depth for each action is up to 
the user to define.


## Usage






 

## Infrastructure as Configuration

## Prototyping

### ddo as container

`docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v $(pwd):/rr -v ~/.azure:/root/.azure ttnesby/ddo:latest ce navutv cr`

### ddo as binary

`./build/darwin/arm64/ddo -no-result ce`



## Test runs

> Pre-requisites: 
- `linux` or `macOS` with `go`,`docker`,`az cli` installed
- `az login` to authenticate with azure has been run

### Export config (ce)

- configs across tenants
`go run ./cmd/ddo/ ./test/infrastructure/automation ce`

- configs for a specific tenant navutv|navno
`go run ./cmd/ddo/ ./test/infrastructure/automation ce navutv`

- config for a specific component rg|cr
`go run ./cmd/ddo/ ./test/infrastructure/automation ce navutv rg`

### Validate config (va)

- configs across tenants
`go run ./cmd/ddo/ ./test/infrastructure/automation va`

- configs for a specific tenant navutv|navno
`go run ./cmd/ddo/ ./test/infrastructure/automation va navutv`

- config for a specific component rg|cr
`go run ./cmd/ddo/ ./test/infrastructure/automation va navutv rg`

**Observes the following:** validation will fail for `cr` if the `rg` is not present in azure

### What-if config (if)

as above, but with `if` instead of `va`

### Deploy config (de)

as above, but with `de` instead of `va`

## Clear cache once a while:

build cache: `go clean -cache`

test cache: `go clean -testcache`

## Verbose tests:

`go test -v ./... -tags=unit -count 10`

`go test -v ./... -tags=integration -count 2`

## Docker stuff

In case you want to run the docker image locally, you can do so by running the following command
`docker run --platform=linux/amd64 -v $(pwd):/rr -it --rm ttnesby/azbicue:latest`