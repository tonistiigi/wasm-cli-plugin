package commands

import (
	"os"
	"path/filepath"

	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tonistiigi/wasm-cli-plugin/control"
)

func NewRootCmd(name string, isPlugin bool, dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Run wasm containers",
		Use:   name,
	}

	var debug bool
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if isPlugin {
			return plugin.PersistentPreRunE(cmd, args)
		}
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	var opt control.Opt

	flags := cmd.PersistentFlags()
	flags.StringVar(&opt.Root, "data-root", filepath.Join(homeDir, ".docker/wasm"), "Root directory of persistent state")
	flags.BoolVarP(&debug, "debug", "D", false, "Enable debug logs")

	addCommands(cmd, dockerCli, &opt)
	return cmd
}

func addCommands(cmd *cobra.Command, dockerCli command.Cli, opt *control.Opt) {
	cmd.AddCommand(
		pullCmd(dockerCli, *opt),
		rmCmd(dockerCli, *opt),
		lsCmd(dockerCli, *opt),
		runCmd(dockerCli, *opt),
	)
}

func getController(opt control.Opt) (*control.Controller, error) {
	if opt.Root == "" {
		opt.Root = "./.state"
	}
	return control.New(opt)
}
