package control

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/platforms"
	bkidentity "github.com/moby/buildkit/identity"
	"github.com/opencontainers/image-spec/identity"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	copy "github.com/tonistiigi/fsutil/copy"
	"github.com/tonistiigi/wasm-cli-plugin/util/singlemounter"
)

func (c *Controller) Run(ctx context.Context, img *images.Image, platform platforms.MatchComparer) error {
	ctx = addNS(ctx)

	chain, err := img.RootFS(ctx, c.cs, platform)
	if err != nil {
		return err
	}

	config, err := img.Config(ctx, c.cs, platform)
	if err != nil {
		return err
	}

	chainID := identity.ChainID(chain)

	sn := c.mdb.Snapshotter("native")

	id := bkidentity.NewID()

	mounts, err := sn.View(ctx, id, chainID.String())
	if err != nil {
		return err
	}

	lm := singlemounter.SingleMounter(mounts)
	mp, err := lm.Mount()
	if err != nil {
		return err
	}

	defer func() {
		if lm != nil {
			lm.Unmount()
		}
	}()

	target := filepath.Join(c.root, "containers", id)
	if err := os.MkdirAll(target, 0700); err != nil {
		return err
	}

	defer func() {
		os.RemoveAll(target)
	}()

	if err := copy.Copy(ctx, mp, ".", target, "."); err != nil {
		return err
	}

	dt, err := content.ReadBlob(ctx, c.cs, config)
	if err != nil {
		return err
	}

	var ociimg ocispec.Image
	err = json.Unmarshal(dt, &ociimg)

	args := append(ociimg.Config.Entrypoint, ociimg.Config.Cmd...)

	args[0] = filepath.Join(target, args[0]) // TODO: not safe

	newArgs := []string{}
	for _, v := range ociimg.Config.Env {
		newArgs = append(newArgs, "--env="+v)
	}
	newArgs = append(newArgs, "--mapdir=/:"+target)
	args = append(newArgs, args...)

	cmd := exec.Command("wasmtime", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
