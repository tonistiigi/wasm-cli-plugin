package commands

import (
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
	"github.com/tonistiigi/wasm-cli-plugin/control"
)

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

	var opt control.Opt

	flags := cmd.PersistentFlags()
	flags.StringVar(&opt.Root, "data-root", "~/.docker/wasm", "Root directory of persistent state")

	addCommands(cmd, dockerCli, &opt)
	return cmd
}

func addCommands(cmd *cobra.Command, dockerCli command.Cli, opt *control.Opt) {
	cmd.AddCommand(
		pullCmd(dockerCli, *opt),
		rmCmd(dockerCli, *opt),
		lsCmd(dockerCli, *opt),
	)
}

func getController(opt control.Opt) (*control.Controller, error) {
	if opt.Root == "" {
		opt.Root = "./.state"
	}
	return control.New(opt)
}
