package main

import (
	"runtime"
)

// Build Goreleaser
func (g *Goreleaser) Build(
	// Target OS to build
	// +default="linux"
	os string,
	// Target architecture to build
	// +optional
	arch string,
) *File {
	if arch == "" {
		arch = runtime.GOARCH
	}
	return g.BuildEnv().
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch).
		WithExec([]string{"go", "build", "-o", "/src/dist/goreleaser"}).
		File("/src/dist/goreleaser")
}

// Base container to build Goreleaser
func (g *Goreleaser) Base() *Container {
	// Base image with Go
	env := dag.Container().
		From("cgr.dev/chainguard/wolfi-base").
		WithExec([]string{"apk", "add", "go"})

	// Mount the Go cache
	env = env.
		WithMountedCache(
			"/go",
			dag.CacheVolume("goreleaser-goroot"),
			ContainerWithMountedCacheOpts{
				Owner: "nonroot",
			}).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod")

	// Mount the Go build cache
	env = env.
		WithMountedCache(
			"/gocache",
			dag.CacheVolume("goreleaser-gobuild"),
			ContainerWithMountedCacheOpts{
				Owner: "nonroot",
			}).
		WithEnvVariable("GOCACHE", "/gocache")

	// Mount the source code
	env = env.
		WithMountedDirectory("/src", g.Source, ContainerWithMountedDirectoryOpts{
			Owner: "nonroot",
		}).
		WithWorkdir("/src")

	return env
}

// Container to build Goreleaser
func (g *Goreleaser) BuildEnv() *Container {
	// Base image with Go
	env := g.Base()

	// Mount the source code
	env = env.With(WithSource(g))

	return env
}

func WithSource(g *Goreleaser) WithContainerFunc {
	return func(c *Container) *Container {
		return c.
			WithMountedDirectory("/src", g.Source, ContainerWithMountedDirectoryOpts{
				Owner: "nonroot",
			}).
			WithWorkdir("/src")
	}
}
