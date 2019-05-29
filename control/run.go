package control

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/platforms"
	bkidentity "github.com/moby/buildkit/identity"
	"github.com/opencontainers/image-spec/identity"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	copy "github.com/tonistiigi/fsutil/copy"
	"github.com/tonistiigi/wasm-cli-plugin/util/singlemounter"
)

type ProcessOpt struct {
	Args       []string
	Entrypoint string
	Env        map[string]string
	Volumes    map[string]string
	Runtime    string
}

func (c *Controller) Run(ctx context.Context, img *images.Image, platform platforms.MatchComparer, po ProcessOpt) error {
	if po.Runtime == "" {
		rt, err := detectRuntime()
		if err != nil {
			return err
		}
		po.Runtime = rt
	}

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
	if err := json.Unmarshal(dt, &ociimg); err != nil {
		return err
	}

	if po.Entrypoint != "" {
		ociimg.Config.Entrypoint = []string{po.Entrypoint}
	}

	if po.Entrypoint != "" || len(po.Args) > 0 {
		ociimg.Config.Cmd = po.Args
	}

	for k, v := range po.Env {
		ociimg.Config.Env = append(ociimg.Config.Env, k+"="+v)
	}

	args := append(ociimg.Config.Entrypoint, ociimg.Config.Cmd...)

	args[0] = filepath.Join(target, args[0]) // TODO: not safe

	switch po.Runtime {
	case "wasmtime":
		newArgs := []string{}
		for _, v := range ociimg.Config.Env {
			parts := strings.SplitN(v, "=", 2)
			if _, ok := po.Env[parts[0]]; !ok {
				newArgs = append(newArgs, "--env="+v)
			}
		}
		for k, v := range po.Env {
			newArgs = append(newArgs, "--env="+k+"="+v)
		}

		newArgs = append(newArgs, "--mapdir=/:"+target)

		for src, dest := range po.Volumes {
			newArgs = append(newArgs, "--mapdir="+dest+":"+src)
		}
		newArgs = append(newArgs, args[0], "--")
		args = append(newArgs, args[1:]...)
	case "wasmer":
		newArgs := []string{"run", "--mapdir=/:" + makeRel(target), args[0], "--"}
		args = append(newArgs, args[1:]...)
	default:
		return errors.Errorf("unknown runtime %s", po.Runtime)
	}

	logrus.Debugf("running: %s %s", po.Runtime, strings.Join(args, " "))

	cmd := exec.Command(binary(po.Runtime), args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to run %s %s", po.Runtime, strings.Join(args, " "))
	}
	return nil
}

func detectRuntime() (string, error) {
	for _, test := range []string{"wasmtime", "wasmer"} {
		if _, err := exec.LookPath(binary(test)); err == nil {
			return test, nil
		}
	}
	return "", errors.Errorf("failed to find and wasm runtimes (wasmtime, wasmer)")
}

func binary(in string) string {
	if runtime.GOOS == "windows" {
		return in + ".exe"
	}
	return in
}

func makeRel(p string) string {
	base, err := os.Getwd()
	if err != nil {
		return p
	}
	p, err = filepath.Rel(base, p)
	if err != nil {
		return p
	}
	if p == "" {
		return p
	}
	return filepath.ToSlash(p)
}
