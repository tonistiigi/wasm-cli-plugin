package main

import (
	"fmt"
	"os"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
	"github.com/tonistiigi/wasm-cli-plugin/commands"
	"github.com/tonistiigi/wasm-cli-plugin/version"
)

func main() {
	if os.Getenv("DOCKER_CLI_PLUGIN_ORIGINAL_CLI_COMMAND") == "" {
		if len(os.Args) < 2 || os.Args[1] != manager.MetadataSubcommandName {
			dockerCli, err := command.NewDockerCli()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			rootCmd := commands.NewRootCmd("docker-wasm", false, dockerCli)
			if err := rootCmd.Execute(); err != nil {
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		return commands.NewRootCmd("wasm", true, dockerCli)
	},
		manager.Metadata{
			SchemaVersion: "0.1.0",
			Vendor:        "Docker Inc.",
			Version:       version.Version,
		})
}
