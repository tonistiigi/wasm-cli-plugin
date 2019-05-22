package commands

import (
	"fmt"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func runRm(dockerCli command.Cli) error {
	// ctx := appcontext.Context()
	return nil
}

func rmCmd(dockerCli command.Cli, opt *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm NAME",
		Short: "Remove an image",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("rm %+v %+v\n", args, opt)
			return runRm(dockerCli)
		},
	}

	return cmd
}
