package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	if err := build(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// define pipeBuild matrix
	oses := []string{"linux", "darwin"}
	arches := []string{"arm64"}

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer func(client *dagger.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	// get reference to the local project
	src := client.Host().Directory(".")

	// create empty directory to put pipeBuild outputs
	outputs := client.Directory()

	// get `golang` image
	golang := client.Container().From("golang:latest")

	// mount cloned repository into `golang` image
	golang = golang.WithMountedDirectory("/rr", src).WithWorkdir("/rr/cmd/ddo")

	for _, goos := range oses {
		for _, goarch := range arches {
			// create a directory for each os and arch
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)

			// set GOARCH and GOOS in the pipeBuild environment
			build := golang.WithEnvVariable("GOOS", goos)
			build = build.WithEnvVariable("GOARCH", goarch)

			// pipeBuild application
			build = build.WithExec([]string{"go", "build", "-o", path})

			// get reference to pipeBuild output directory in container
			outputs = outputs.WithDirectory(path, build.Directory(path))
		}
	}
	// write pipeBuild artifacts to host
	_, err = outputs.Export(ctx, ".")
	if err != nil {
		return err
	}

	return nil
}
