package commands

import (
	"github.com/containerd/containerd/platforms"
	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/moby/buildkit/util/appcontext"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tonistiigi/wasm-cli-plugin/control"
)

func runPull(dockerCli command.Cli, opt control.Opt, ref string) error {
	ctx := appcontext.Context()

	c, err := getController(opt)
	if err != nil {
		return err
	}
	defer c.Close()

	img, err := c.Pull(ctx, ref, allowedPlatforms())
	if err != nil {
		return errors.WithStack(err)
	}

	_ = img

	return nil
}

func pullCmd(dockerCli command.Cli, opt control.Opt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull REF",
		Short: "Pull a wasm image",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPull(dockerCli, opt, args[0])
		},
	}

	return cmd
}

func allowedPlatforms() platforms.MatchComparer {
	return platforms.Only(ocispec.Platform{
		OS:           "wasi",
		Architecture: "wasm",
	})
}
