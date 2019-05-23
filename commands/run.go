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

func runRun(dockerCli command.Cli, opt control.Opt, ref string) error {
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

	pm := platforms.Only(p)

	img, err := c.GetRuntimeImage(ctx, ref, pm)
	if err != nil || img == nil {
		img, err = c.Pull(ctx, ref, pm)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if err := c.Run(ctx, img, pm); err != nil {
		return err
	}

	return nil
}

func runCmd(dockerCli command.Cli, opt control.Opt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run REF",
		Short: "Run an image",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRun(dockerCli, opt, args[0])
		},
	}

	return cmd
}
