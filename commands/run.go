package commands

import (
	"os"
	"strings"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tonistiigi/wasm-cli-plugin/control"
)

type runOpt struct {
	entrypoint string
	args       []string
	volumes    []string
	env        []string
	runtime    string
}

func runRun(dockerCli command.Cli, opt control.Opt, ref string, ro runOpt) error {
	ctx := appcontext.Context()

	c, err := getController(opt)
	if err != nil {
		return err
	}
	defer c.Close()

	pm := allowedPlatforms()

	img, err := c.GetRuntimeImage(ctx, ref, pm)
	if err != nil || img == nil {
		img, err = c.Pull(ctx, ref, pm)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	po := control.ProcessOpt{
		Args:       ro.args,
		Entrypoint: ro.entrypoint,
		Runtime:    ro.runtime,
	}

	if len(ro.env) > 0 {
		m := map[string]string{}
		for _, env := range ro.env {
			parts := strings.SplitN(env, "=", 2)
			var v string
			if len(parts) == 2 {
				v = parts[1]
			} else {
				v = os.Getenv(parts[0])
			}
			m[parts[0]] = v
		}
		po.Env = m
	}

	if len(ro.volumes) > 0 {
		m := map[string]string{}
		for _, v := range ro.volumes {
			parts := strings.SplitN(v, ":", 2)
			if len(parts) != 2 {
				return errors.Errorf("invalid volume %q, only bind mounts supported", v)
			}
			m[parts[0]] = parts[1]
		}
		po.Volumes = m
	}

	if err := c.Run(ctx, img, pm, po); err != nil {
		return err
	}

	return nil
}

func runCmd(dockerCli command.Cli, opt control.Opt) *cobra.Command {
	var ro runOpt

	cmd := &cobra.Command{
		Use:   "run REF",
		Short: "Run a wasm image",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ro.args = args[1:]
			return runRun(dockerCli, opt, args[0], ro)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&ro.entrypoint, "entrypoint", "", "Overwrite the default ENTRYPOINT of the image")
	flags.StringSliceVarP(&ro.env, "env", "e", nil, "Set environment variables")
	flags.StringSliceVarP(&ro.volumes, "volume", "v", nil, "Bind mount a volume")
	flags.StringVar(&ro.runtime, "runtime", "", "WASM runtime")

	return cmd
}
