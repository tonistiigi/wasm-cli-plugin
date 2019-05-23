package commands

import (
	"fmt"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/spf13/cobra"
	"github.com/tonistiigi/wasm-cli-plugin/control"
)

func runLs(dockerCli command.Cli, opt control.Opt) error {
	ctx := appcontext.Context()

	c, err := getController(opt)
	if err != nil {
		return err
	}
	defer c.Close()

	imgs, err := c.Images(ctx)
	if err != nil {
		return err
	}

	for _, img := range imgs {
		fmt.Printf("img: %+v\n", img)
	}

	return nil
}

func lsCmd(dockerCli command.Cli, opt control.Opt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List images",
		Args:  cli.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLs(dockerCli, opt)
		},
	}

	return cmd
}
