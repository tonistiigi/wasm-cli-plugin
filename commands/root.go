package commands

import (
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

type rootOptions struct {
	root string
}

func NewRootCmd(name string, isPlugin bool, dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Run wasm containers",
		Use:   name,
	}
	if isPlugin {
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			return plugin.PersistentPreRunE(cmd, args)
		}
	}

	var options rootOptions

	flags := cmd.PersistentFlags()
	flags.StringVar(&options.root, "data-root", "~/.docker/wasm", "Root directory of persistent state")

	addCommands(cmd, dockerCli, &options)
	return cmd
}

func addCommands(cmd *cobra.Command, dockerCli command.Cli, opt *rootOptions) {
	cmd.AddCommand(
		pullCmd(dockerCli, opt),
		rmCmd(dockerCli, opt),
		lsCmd(dockerCli, opt),
	)
}
