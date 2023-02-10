# ddo

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