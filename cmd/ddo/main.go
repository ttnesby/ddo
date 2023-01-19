package main

import (
	"context"
	"dagger.io/dagger"
	"ddo/path"
	"fmt"
)

func main() {
	if err := build(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func build(ctx context.Context) error {

	fmt.Println("Start dagger and initialize client")
	client, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}
	defer func(client *dagger.Client) {
		_ = client.Close()
	}(client)

	fmt.Printf("Connect to host repository %s\n", path.RepoRoot())
	repo := client.Host().Directory(path.RepoRoot())

	da := path.HomeAbs(".azure")
	fmt.Printf("Connect to host %s\n", da)
	if !path.AbsExists(da) {
		return fmt.Errorf("folder %s does not exist", da)
	}

	dotAzure := client.Host().Directory(
		da,
		dagger.HostDirectoryOpts{
			Include: []string{
				"azureProfile.json",
				"msal_http_cache.bin",
				"msal_token_cache.json",
				"service_principal_entries.json",
			},
		},
	)

	contRef := "docker.io/ttnesby/azbicue:latest"
	fmt.Printf("Connect to container %s\n", contRef)

	azbicue := client.Container().
		From(contRef).
		WithMountedDirectory("rr", repo).
		WithMountedDirectory("/root/.azure", dotAzure).
		WithWorkdir("/rr")

	actionsCmd := []string{
		"cue",
		"export",
		"-p",
		"ddospec",
		"./test/ddo.cue",
		"./test/ddo.schema.cue",
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
