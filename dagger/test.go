package main

import (
	"context"
	"fmt"
	"runtime"
)

// Test Goreleaser
func (g *Goreleaser) Test(ctx context.Context) (string, error) {
	return g.TestEnv().
		WithExec([]string{
			"go",
			"test",
			"-failfast",
			// "-race", // TODO: change base
			"-coverpkg=./...",
			"-covermode=atomic",
			"-coverprofile=coverage.txt",
			"./...",
			"-run",
			".",
		}).
		Stdout(ctx)
}

// Container to test Goreleaser
func (g *Goreleaser) TestEnv() *Container {
	// Dependencies needed for testing
	testDeps := []string{
		"bash",
		"curl",
		"git",
		"gpg",
		"gpg-agent",
		// "nix",
		"upx",
		"cosign",
		"docker",
		"syft",
	}
	return g.BuildEnv().
		// WithEnvVariable("CGO_ENABLED", "1"). // TODO: change base
		WithServiceBinding("localhost", dag.Docker().Engine()). // TODO: fix localhost
		WithEnvVariable("DOCKER_HOST", "tcp://localhost:2375").
		WithExec(append([]string{"apk", "add"}, testDeps...)).
		With(installNix).
		With(installBuildx).
		// WithExec([]string{"sh", "-c", "sh <(curl -L https://nixos.org/nix/install) --no-daemon"})
		// WithExec([]string{"chown", "-R", "nonroot", "/nix"}).
		WithUser("nonroot").
		WithExec([]string{"go", "install", "github.com/google/ko@latest"})
}

func installNix(target *Container) *Container {
	nix := dag.Container().From("nixos/nix")
	nixBin := "/root/.nix-profile/bin"

	binaries := []string{
		"nix",
		"nix-build",
		"nix-channel",
		"nix-collect-garbage",
		"nix-copy-closure",
		"nix-daemon",
		"nix-env",
		"nix-hash",
		"nix-instantiate",
		"nix-prefetch-url",
		"nix-shell",
		"nix-store",
	}

	for _, binary := range binaries {
		target = target.WithFile("/bin/"+binary, nix.File(nixBin+"/"+binary))
	}

	target = target.WithDirectory("/nix/store", nix.Directory("/nix/store"))

	return target
}

func installBuildx(target *Container) *Container {
	arch := runtime.GOARCH
	url := fmt.Sprintf("https://github.com/docker/buildx/releases/download/v0.15.1/buildx-v0.15.1.linux-%s", arch)

	bin := dag.HTTP(url)

	return target.WithFile(
		"/usr/lib/docker/cli-plugins/docker-buildx",
		bin,
		ContainerWithFileOpts{
			Permissions: 0777,
		})
}
