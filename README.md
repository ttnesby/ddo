# ddo

## Clear cache once a while:

build cache: `go clean -cache`

test cache: `go clean -testcache`

## Verbose tests:

`go test -v ./... -tags=unit -count 10`

`go test -v ./... -tags=integration -count 2`
