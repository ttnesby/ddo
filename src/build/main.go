package main

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"os"
)

func main() {

	exitCode := func() int {
		if err := build(context.Background()); err != nil {
			fmt.Println(err)
			return 1
		}
		return 0
	}()
	os.Exit(exitCode)
}

func build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// define pipeBuild matrix
	oses := []string{"darwin"}  //[]string{"linux", "darwin"}
	arches := []string{"arm64"} //[]string{"arm64", "amd64"}

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

	// create empty directory to put pipeBuild hostFolder
	hostFolder := client.Directory()

	// get `golang` image
	cnt := client.
		Container().From("golang:latest").
		WithMountedDirectory(
			"/rr",
			client.Host().Directory(
				".",
				dagger.HostDirectoryOpts{Exclude: []string{"build/"}},
			),
		).
		WithWorkdir("/rr/cmd/ddo")

	for _, goos := range oses {
		for _, goarch := range arches {
			// create a directory for each os and arch
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)

			// set GOARCH and GOOS in the pipeBuild environment
			cnt = cnt.WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithExec([]string{"go", "build", "-o", path})

			// get reference to pipeBuild output directory in container
			hostFolder = hostFolder.WithDirectory(path, cnt.Directory(path))
		}
	}
	// write pipeBuild artifacts to host
	_, err = hostFolder.Export(ctx, ".")
	if err != nil {
		return err
	}

	return nil
}
