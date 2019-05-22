package commands

import (
	"fmt"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func runPull(dockerCli command.Cli) error {
	// ctx := appcontext.Context()
	return nil
}

func pullCmd(dockerCli command.Cli, opt *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull REF",
		Short: "Pull an image",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("pull %+v %+v\n", args, opt)
			return runPull(dockerCli)
		},
	}

	return cmd
}
