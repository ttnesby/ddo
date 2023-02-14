# Infrastructure as Configuration

`ddo` (DaggerDO) is a tool for `azure infrastructure as configuration`. The building blocks are:
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
terms of just one template, the what-if is working as good-as-it-can. The `false-positives`, depending on the type of 
infrastructure component, are still popping up.

> `ddo` is `path oriented`, meaning that the chosen action will be applied to all components in the path. See 
> [Orchestration configuration](#orchestration-configuration) for more details.

## Configuration example

`./test/infrastructure` is an example of a simple infrastructure solution. A resource group with a container registry, 
to be deployed to two different tenants. 

>**Warning** cuelang versus json/yaml is an involving topic beyond the scope of this document. Briefly, it's a 
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
in `Orchestration configuration` and in GitHub actions, e.g. federated login.

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

The `json-look` with action `ce` only is displayed below. The json for the other actions are similar.
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
#### Path oriented
By looking at the above json-file, the meaning of `path oriented` is easy to understand, it's a kind of path in the json file:
- `ddo ce navutv rg` will export the config for `resourceGroup` with `tenant=navutv` tag
- `ddo ce navutv cr` will export the config for `containerRegistry` with `tenant=navutv` tag
- `ddo ce navno` will export the config for both `resourceGroup` and `containerRegistry` with `tenant=navno` tag
- `ddo ce` will export the config for both components and both tenants

The relevant structure and depth for each action is up to the user to define.

#### deployOrder and parallelism

The `deployOrder` is a list of lists, defining the order of component deployment. All components in a sub list are invoked
in parallel. See order and parallelism per action:
> - `ce`,`va`,`if` - order doesn't matter, all components are in one sublist
> - `de` - order matters, computation will stop if a component fails
> - `evomer` - reversed order matters, computation will stop if a component fails 

### Cyclic dependencies

Depending on the infrastructure solution, there might be cyclic dependencies between components. E.g. in a solution for 
`hybrid`connectivity, the virtual hub is created first, while azure firewall and express route gateway is created 
later, both importing the virtual hub configuration and `link` to the virtual hub. 

If the virtual hub is updated (tags or whatever) later, the hub must refer to the resource Ids for firewall and 
express route gateway iff they exist. Since cue configuration is `hermetic` and cyclic dependencies are not allowed, such 
data dependencies are solved with use of tags.

The `resourceGroup` component is over-engineered, but with the purpose of showing how to solve cyclic dependencies.

The `resourceGroup/deploy.cue` component is defined with the following extra tags:
```cue
_crid:             *"" | string  @tag(crid)
_crpoltruststatus: *"" | string  @tag(crpoltruststatus)
_crb64:            *"" | string  @tag(crb64)
```
They are defaulting to empty string, but will be set by the orchestration, when the `resourceGroup` component is engaged.
In this example, the extra tags are just enhancing the resource group tags.

The value for those tags are defined the un orchestration configuration (`ddo.cue`):
```cue
rg: #component & {
     folder: "\(_componentsPath)/resourceGroup"
     tags: [
        "tenant=\(#tenant)",
        // azure data lookup, << pathComponent << pathData (`azure resource show` as reference)
        "crid=<<\(#tenant).cr<<id",
        "crpoltruststatus=<<\(#tenant).cr<<properties.policies.trustPolicy.status",
        // azure data lookup, << pathComponent << b64
        "crb64=<<\(#tenant).cr<<b64",
    ]
}
```

The syntax is simple; `tagname=<<componentPath<<dataPath|b64`. 
- `crid=<<\(#tenant).cr<<id` will give the `id` for container registry in relevant tenant
- `crpoltruststatus=<<\(#tenant).cr<<properties.policies.trustPolicy.status` will give the `properties.policies.trustPolicy.status` 
for container registry in relevant tenant
- `crb64=<<\(#tenant).cr<<b64` will give the `b64` string for container registry properties in relevant tenant

The `b64` is a special case since tags are simple values only, no objects or lists.

> The data set reference for a component is `azure resource show --ids <#resourceId of the component>`

## Usage

[`az cli`](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli) must be installed. 
`ddo` supports inheritance of the host `.azure` folder

> - Pre-requisites: a terminal window and in a repo root folder supporting relevant configuration, e.g. this project
> - Pre-requisites: in a terminal window, do az login for relevant tenants 
> - Do `docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v $(pwd):/rr -v ~/.azure:/root/.azure docker.io/ttnesby/ddo:latest`

The docker run command has three mounts:
1. `/var/run/docker.sock:/var/run/docker.sock` - to be able to run docker commands from within the container (kind of docker-in-docker where dagger.io will be a `sibling` of this container)
2. `$(pwd):/rr` - to be able to access the current folder from within the container
3. `~/.azure:/root/.azure` - to be able to access the hosts `.azure` folder from within the container

From now on, `docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v $(pwd):/rr -v ~/.azure:/root/.azure docker.io/ttnesby/ddo:latest` 
is abbreviated to `./output/darwin/arm64/ddo` (the binary version).

### Help
```zsh
./output/darwin/arm64/ddo -h
#output
usage: ./output/darwin/arm64/ddo
ddo [options] [action path...]

Action path is a path to component in ddo.cue file, starting with one of:
ce - config export
va - validate config against azure
if - what-if analysis against azure
de - deploy to azure
evomer - remove iff the component has '#resourceId' definition

e.g.
ddo ce navutv rg - config export of navutv and component rg
ddo ce navutv - config export of all components in navutv
ddo if - what-if of all components in all tenants

Options:
  -cnt string
    	container ref. hosting az cli, bicep and cue. Default is 'docker.io/ttnesby/azbicue:latest' (default "docker.io/ttnesby/azbicue:latest")
  -debug
    	debug mode
  -debug-container
    	debug mode for dagger.io
  -no-result
    	No display of action result
```
> - `-no-result` is a good thing to use when running in a CI/CD pipeline. The exit code will be 0 if the actions were 
successful, otherwise 1.
> - `-debug` will show the orchestration of `cue cli` and `az cli` commands
> - `-debug-container` will show the details inside the `sibling` container used by dagger.io

### -no-result ce navutv 
```zsh
./output/darwin/arm64/ddo -no-result ce navutv
#output
2023-02-13T14:00:12+01:00 INF action.go:43 > start dagger client
2023-02-13T14:00:14+01:00 INF action.go:155 > searched for ddo.cue [./test/infrastructure/automation]
2023-02-13T14:00:14+01:00 INF action.go:160 > reading action specification ./test/infrastructure/automation
2023-02-13T14:00:15+01:00 INF action.go:190 > get selection: actions.ce.navutv|@pretty
2023-02-13T14:00:15+01:00 INF action.go:190 > get selection: deployOrder|@pretty
2023-02-13T14:00:15+01:00 INF action.go:141 > resolve data injections for components
2023-02-13T14:00:15+01:00 INF component.go:220 > injection for ./test/infrastructure/resourceGroup
2023-02-13T14:00:15+01:00 INF action.go:190 > get selection: actions.ce.navutv.cr|@pretty
2023-02-13T14:00:23+01:00 INF action.go:190 > get selection: actions.ce.navutv.cr|@pretty
2023-02-13T14:00:23+01:00 INF action.go:190 > get selection: actions.ce.navutv.cr|@pretty
2023-02-13T14:00:24+01:00 INF component.go:220 > injection for ./test/infrastructure/containerRegistry
2023-02-13T14:00:24+01:00 INF component.go:171 > [navutv cr] ce
2023-02-13T14:00:24+01:00 INF component.go:171 > [navutv rg] ce
2023-02-13T14:00:25+01:00 INF component.go:110 > [navutv rg] done
2023-02-13T14:00:25+01:00 INF component.go:110 > [navutv cr] done
2023-02-13T14:00:25+01:00 INF action.go:94 > done!
```
### -no-result va navutv
```zsh
./output/darwin/arm64/ddo -no-result va navutv
#output
2023-02-13T14:04:26+01:00 INF action.go:43 > start dagger client
2023-02-13T14:04:28+01:00 INF action.go:155 > searched for ddo.cue [./test/infrastructure/automation]
2023-02-13T14:04:28+01:00 INF action.go:160 > reading action specification ./test/infrastructure/automation
2023-02-13T14:04:30+01:00 INF action.go:190 > get selection: actions.va.navutv|@pretty
2023-02-13T14:04:30+01:00 INF action.go:190 > get selection: deployOrder|@pretty
2023-02-13T14:04:30+01:00 INF action.go:141 > resolve data injections for components
2023-02-13T14:04:30+01:00 INF component.go:220 > injection for ./test/infrastructure/resourceGroup
2023-02-13T14:04:30+01:00 INF action.go:190 > get selection: actions.ce.navutv.cr|@pretty
2023-02-13T14:04:37+01:00 INF action.go:190 > get selection: actions.ce.navutv.cr|@pretty
2023-02-13T14:04:38+01:00 INF action.go:190 > get selection: actions.ce.navutv.cr|@pretty
2023-02-13T14:04:38+01:00 INF component.go:220 > injection for ./test/infrastructure/containerRegistry
2023-02-13T14:04:38+01:00 INF component.go:171 > [navutv cr] va
2023-02-13T14:04:38+01:00 INF component.go:171 > [navutv rg] va
2023-02-13T14:05:06+01:00 INF component.go:110 > [navutv rg] done
2023-02-13T14:05:06+01:00 ERR component.go:108 > error="[navutv cr] failed \ninput:1: container.from.withMountedDirectory.withMountedDirectory.withWorkdir.withExec.withExec.stdout process \"az deployment group validate --name 46e12143-6c37-587d-96c8-15472baf8e89 --subscription ca1e4592-6211-4c03-aac5-e681c1d1ea0d --resource-group container-registry --template-file ./test/infrastructure/containerRegistry/main.bicep --parameters @/tmp/ddo.parameters.01GS5E16GG1JV1BZZN9K362MP8.json --out yaml\" did not complete successfully: exit code: 1\nStdout:\n\nStderr:\nERROR: {\"code\": \"ResourceGroupNotFound\", \"message\": \"Resource group 'container-registry' could not be found.\"}\n\nPlease visit https://dagger.io/help#go for troubleshooting guidance."
2023-02-13T14:05:07+01:00 ERR action.go:250 > error="1 component(s) failed"
```
The reason for cr failure is `az cli` by design. Cannot validate a component to be placed in a resource group if the 
resource group doesn't exist. By deploying the resource group first, both validation and what-if will pass. 

> The user can design composite bicep templates, but then the what-if will be less trustworthy. 

### if navutv rg 
```zsh
./output/darwin/arm64/ddo if navutv rg
#output
2023-02-13T14:11:05+01:00 INF action.go:43 > start dagger client
2023-02-13T14:11:07+01:00 INF action.go:155 > searched for ddo.cue [./test/infrastructure/automation]
2023-02-13T14:11:07+01:00 INF action.go:160 > reading action specification ./test/infrastructure/automation
2023-02-13T14:11:10+01:00 INF action.go:190 > get selection: actions.if.navutv.rg|@pretty
2023-02-13T14:11:10+01:00 INF action.go:190 > get selection: deployOrder|@pretty
2023-02-13T14:11:10+01:00 INF action.go:141 > resolve data injections for components
2023-02-13T14:11:10+01:00 INF component.go:220 > injection for ./test/infrastructure/resourceGroup
2023-02-13T14:11:10+01:00 INF action.go:190 > get selection: actions.ce.navutv.cr|@pretty
2023-02-13T14:11:17+01:00 INF action.go:190 > get selection: actions.ce.navutv.cr|@pretty
2023-02-13T14:11:18+01:00 INF action.go:190 > get selection: actions.ce.navutv.cr|@pretty
2023-02-13T14:11:18+01:00 INF component.go:171 > [navutv rg] if
2023-02-13T14:11:59+01:00 INF component.go:112 > [navutv rg]
Note: The result may contain false positive predictions (noise).
You can help us improve the accuracy of the result by opening an issue here: https://aka.ms/WhatIfIssues

Resource and property changes are indicated with this symbol:
  + Create

The deployment will update the following scope:

Scope: /subscriptions/ca1e4592-6211-4c03-aac5-e681c1d1ea0d

  + resourceGroups/container-registry [2021-04-01]

      apiVersion:           "2021-04-01"
      id:                   "/subscriptions/ca1e4592-6211-4c03-aac5-e681c1d1ea0d/resourceGroups/container-registry"
      location:             "norwayeast"
      name:                 "container-registry"
      tags.Application:     "ddo"
      tags.Provenance:      "https://github.com/ttnesby/ddo"
      tags.Team:            "azure-platform"
      tags.Technical Owner: "azure-platform"
      type:                 "Microsoft.Resources/resourceGroups"

Resource changes: 1 to create.
2023-02-13T14:11:59+01:00 INF action.go:94 > done!
```
### de navutv rg
```zsh
./output/darwin/arm64/ddo de navutv rg
#output
2023-02-13T15:21:47+01:00 INF action.go:43 > start dagger client
2023-02-13T15:21:49+01:00 INF action.go:147 > searched for ddo.cue [./test/infrastructure/automation]
2023-02-13T15:21:49+01:00 INF action.go:152 > reading action specification ./test/infrastructure/automation
2023-02-13T15:21:50+01:00 INF action.go:182 > get selection: actions.de.navutv.rg|@pretty
2023-02-13T15:21:50+01:00 INF action.go:182 > get selection: deployOrder|@pretty
2023-02-13T15:21:50+01:00 INF action.go:133 > resolve data injections for components
2023-02-13T15:21:50+01:00 INF component.go:220 > injection for ./test/infrastructure/resourceGroup
2023-02-13T15:21:50+01:00 INF action.go:182 > get selection: actions.ce.navutv.cr|@pretty
2023-02-13T15:21:56+01:00 INF action.go:182 > get selection: actions.ce.navutv.cr|@pretty
2023-02-13T15:21:56+01:00 INF action.go:182 > get selection: actions.ce.navutv.cr|@pretty
2023-02-13T15:21:57+01:00 INF component.go:171 > [navutv rg] de
2023-02-13T15:22:53+01:00 INF component.go:112 > [navutv rg]
id: /subscriptions/ca1e4592-6211-4c03-aac5-e681c1d1ea0d/providers/Microsoft.Resources/deployments/f4381b42-ff52-586b-b5ee-4fbdcb7f3406
location: norwayeast
name: f4381b42-ff52-586b-b5ee-4fbdcb7f3406
properties:
...shortened...
2023-02-13T15:22:53+01:00 INF action.go:89 > done!
```
Now is it possible to run `va navutv` and `if navutv`.
> `de navutv` and `evomer navutv` is left to the reader.

## Technical details

### ddo cli
- `./cmd/ddo/main.go` using a set of components (action, alogger, arg, azcli, cuecli and path)

### build ddo binary
- `./src/build/main.go` - using dagger.io

### build and push ddo docker image
- `./src/image/main.go` - using dagger.io

### build and push tools docker image (az cli, bicep, cue cli)
- `./src/tools/main.go` - using dagger.io

### cue.mod
- `./cue.mod/module.cue` - is the cue version of go.mod. Used by the cue config. in test folder.
