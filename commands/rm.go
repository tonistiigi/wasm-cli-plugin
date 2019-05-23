package commands

import (
	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/spf13/cobra"
	"github.com/tonistiigi/wasm-cli-plugin/control"
)

func runRm(dockerCli command.Cli, opt control.Opt, ref string) error {
	ctx := appcontext.Context()

	c, err := getController(opt)
	if err != nil {
		return err
	}
	defer c.Close()

	if err := c.Delete(ctx, ref); err != nil {
		return err
	}

	return nil
}

func rmCmd(dockerCli command.Cli, opt control.Opt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm NAME",
		Short: "Remove a wasm image",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRm(dockerCli, opt, args[0])
		},
	}

	return cmd
}
