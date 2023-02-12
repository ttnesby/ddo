package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	platformFormat "github.com/containerd/containerd/platforms"
)

var platforms = []dagger.Platform{
	"linux/amd64",
	"linux/arm64",
}

// the container registry for the multi-platform image
const imageRepo = "docker.io/ttnesby/ddo:latest"

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

// util that returns the architecture of the provided platform
func architectureOf(platform dagger.Platform) string {
	return platformFormat.MustParse(string(platform)).Architecture
}

func build(ctx context.Context) error {

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout), dagger.WithWorkdir("."))
	if err != nil {
		return err
	}
	defer func(client *dagger.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	hostRepoRoot := client.Host().
		Directory(".", dagger.HostDirectoryOpts{Exclude: []string{"build/"}})

	const (
		name       = "/ddo"
		repoRoot   = "/rr"
		outputDir  = "/output"
		sourceCode = repoRoot + "/cmd" + name
		usrBin     = "/usr/bin"
	)

	platformVariants := make([]*dagger.Container, 0, len(platforms))
	for _, platform := range platforms {
		// pull the golang image for the *host platform*. This is
		// accomplished by just not specifying a platform; the default
		// is that of the host.
		ctr := client.Container().
			From("golang:1.19-alpine").
			WithMountedDirectory(repoRoot, hostRepoRoot).
			WithMountedDirectory(outputDir, client.Directory()).
			WithEnvVariable("CGO_ENABLED", "0").
			WithEnvVariable("GOOS", "linux").
			WithEnvVariable("GOARCH", architectureOf(platform)).
			WithWorkdir(sourceCode).
			WithExec([]string{
				"go", "build",
				"-o", outputDir + name,
				sourceCode,
			})

		platformVariants = append(
			platformVariants,
			client.
				Container(dagger.ContainerOpts{Platform: platform}).
				From("alpine:latest").
				// install docker cli
				WithExec([]string{"apk", "update"}).
				WithExec([]string{"apk", "add", "--no-cache", "docker-cli"}).
				// add ddo
				WithFile(usrBin+name, ctr.File(outputDir+name)).
				WithEntrypoint([]string{usrBin + name}),
		)
	}

	// publishing the final image uses the same API as single-platform
	// images, but now additionally specify the `PlatformVariants`
	// option with the containers built before.
	imageDigest, err := client.
		Container().
		Publish(ctx, imageRepo, dagger.ContainerPublishOpts{
			PlatformVariants: platformVariants,
		})
	if err != nil {
		return err
	}
	fmt.Println("published multi-platform image with digest", imageDigest)
	return nil
}
