package control

import (
	"context"

	"github.com/containerd/containerd/namespaces"
	"github.com/docker/distribution/reference"
	"github.com/pkg/errors"
)

func parseRef(rawRef string) (string, error) {
	if rawRef == "" {
		return "", errors.New("missing ref")
	}
	parsed, err := reference.ParseNormalizedNamed(rawRef)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return reference.TagNameOnly(parsed).String(), nil
}

func addNS(ctx context.Context) context.Context {
	return namespaces.WithNamespace(context.Background(), "wasm-cli")
}
