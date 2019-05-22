package commands

import (
	"fmt"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func runLs(dockerCli command.Cli) error {
	// ctx := appcontext.Context()
	return nil
}

func lsCmd(dockerCli command.Cli, opt *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List images",
		Args:  cli.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("images %+v %+v\n", args, opt)
			return runLs(dockerCli)
		},
	}

	return cmd
}
