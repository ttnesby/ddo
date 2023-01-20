package main

import (
	"context"
	"dagger.io/dagger"
	"ddo/path"
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
	fmt.Printf("Exit code: %d\n", exitCode)
	os.Exit(exitCode)
}

func getDotAzurePath() string {
	return path.HomeAbs(".azure")
}

func hostRepoRoot(c *dagger.Client) *dagger.Directory {
	return c.Host().Directory(path.RepoRoot())
}

func hostDotAzure(c *dagger.Client) *dagger.Directory {
	return c.Host().Directory(getDotAzurePath(),
		dagger.HostDirectoryOpts{
			Include: []string{
				"azureProfile.json",
				"msal_http_cache.bin",
				"msal_token_cache.json",
				"service_principal_entries.json",
			},
		},
	)
}

func build(ctx context.Context) error {

	const (
		containerRef      = "docker.io/ttnesby/azbicue:latest"
		containerDotAzure = "/root/.azure"
		containerRepoRoot = "/rr"
	)

	fmt.Println("Start dagger client")
	client, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}

	defer func(client *dagger.Client) {
		_ = client.Close()
	}(client)

	fmt.Printf("Verify and connect to host repository %s\n", path.RepoRoot())
	fmt.Printf("Verify and connect to host %s\n", getDotAzurePath())
	if !path.AbsExists(getDotAzurePath()) {
		return fmt.Errorf("folder %s does not exist", getDotAzurePath())
	}

	fmt.Printf("Start container %s mounting [repo root, .azure]\n", containerRef)

	azbicue := client.Container().
		From(containerRef).
		WithMountedDirectory(containerRepoRoot, hostRepoRoot(client)).
		WithMountedDirectory(containerDotAzure, hostDotAzure(client)).
		WithWorkdir(containerRepoRoot)

	actionsCmd := []string{
		"cue",
		"export",
		"-p",
		"ddospec",
		"./test/ddo.cue",
		"./test/actions.schema.cue",
	}

	fmt.Printf("Reading action specification %v\n\n", actionsCmd)
	actionsJson, err := azbicue.WithExec(actionsCmd).Stdout(ctx)

	if err != nil {
		return err
	}

	fmt.Printf("%v", actionsJson)

	println("Done!")

	return nil
}
