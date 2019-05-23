package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

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

	tw := tabwriter.NewWriter(os.Stdout, 1, 8, 1, '\t', 0)

	fmt.Fprintln(tw, "NAME\tDIGEST")

	for _, img := range imgs {
		fmt.Fprintf(tw, "%s\t%s\n", img.Name, img.Target.Digest)
	}

	tw.Flush()

	return nil
}

func lsCmd(dockerCli command.Cli, opt control.Opt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List wasm images",
		Args:  cli.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLs(dockerCli, opt)
		},
	}

	return cmd
}
