package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

var platforms = []dagger.Platform{
	"linux/amd64",
	"linux/arm64",
}

// the container registry for the multi-platform image
const imageRepo = "docker.io/ttnesby/tooling:latest"

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

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout), dagger.WithWorkdir("."))
	if err != nil {
		return err
	}
	defer func(client *dagger.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	platformVariants := make([]*dagger.Container, 0, len(platforms))
	for _, platform := range platforms {

		ctrAzCliBicep := client.Container(dagger.ContainerOpts{Platform: platform}).
			From("mcr.microsoft.com/azure-cli:latest").
			//https://learn.microsoft.com/en-gb/azure/azure-resource-manager/bicep/install#linux
			WithExec([]string{
				"curl",
				"-Lo",
				"bicep",
				"https://github.com/Azure/bicep/releases/latest/download/bicep-linux-musl-x64",
			}).
			WithExec([]string{
				"chmod",
				"+x",
				"./bicep",
			})

		ctrCue := client.Container(dagger.ContainerOpts{Platform: platform}).
			From("docker.io/cuelang/cue:latest")

		ctrAzCliBiCue := client.Container(dagger.ContainerOpts{Platform: platform}).
			From("mcr.microsoft.com/azure-cli:latest").
			WithFile("/usr/bin/cue", ctrCue.File("/usr/bin/cue")).
			WithFile("/usr/local/bin/bicep", ctrAzCliBicep.File("/bicep"))

		platformVariants = append(
			platformVariants,
			ctrAzCliBiCue,
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
