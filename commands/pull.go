package commands

import (
	"github.com/containerd/containerd/platforms"
	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/moby/buildkit/util/appcontext"
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

	p, err := platforms.Parse("wasi/wasm")
	if err != nil {
		return errors.Wrapf(err, "invalid platform")
	}

	img, err := c.Pull(ctx, ref, platforms.Only(p))
	if err != nil {
		return errors.WithStack(err)
	}

	_ = img

	return nil
}

func pullCmd(dockerCli command.Cli, opt control.Opt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull REF",
		Short: "Pull an image",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPull(dockerCli, opt, args[0])
		},
	}

	return cmd
}
