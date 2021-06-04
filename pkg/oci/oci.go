package oci

import (
	"context"
	"os"

	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/oras-project/oras-go/pkg/content"
	"github.com/oras-project/oras-go/pkg/oras"
	"github.com/sirupsen/logrus"
)

const (
	haulerMediaType = "application/vnd.oci.image"
)

func Get(ctx context.Context, src string, dst string) error {

	store := content.NewFileStore(dst)
	defer store.Close()

	resolver, err := resolver()
	if err != nil {
		return err
	}

	allowedMediaTypes := []string{
		haulerMediaType,
	}

	// Pull file(s) from registry and save to disk
	logrus.Infof("Pulling from %s and saving to %s\n", src, dst)
	desc, _, err := oras.Pull(ctx, resolver, src, store, oras.WithAllowedMediaTypes(allowedMediaTypes))

	if err != nil {
		return err
	}

	logrus.Infof("Pulled from %s with digest %s\n", src, desc.Digest)

	return nil
}

func Put(ctx context.Context, src string, dst string) error {

	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	resolver, err := resolver()
	if err != nil {
		return err
	}

	store := content.NewMemoryStore()

	contents := []ocispec.Descriptor{
		store.Add(src, haulerMediaType, data),
	}

	desc, err := oras.Push(ctx, resolver, dst, store, contents)
	if err != nil {
		return err
	}

	logrus.Infof("pushed %s to %s with digest: %s", src, dst, desc.Digest)

	return nil
}

func resolver() (remotes.Resolver, error) {
	resolver := docker.NewResolver(docker.ResolverOptions{PlainHTTP: true})
	return resolver, nil
}
