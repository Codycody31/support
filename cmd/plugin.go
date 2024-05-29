package cmd

import (
	"go.codycody31.dev/support/plugins"

	"github.com/urfave/cli/v2"
)

var PluginCommand = &cli.Command{
	Name:  "plugin",
	Usage: "Manage plugins",
	Subcommands: []*cli.Command{
		{
			Name:   "enable",
			Usage:  "Enable a plugin",
			Action: plugins.EnablePlugin,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "Plugin name",
					Required: true,
				},
			},
		},
		{
			Name:   "disable",
			Usage:  "Disable a plugin",
			Action: plugins.DisablePlugin,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "Plugin name",
					Required: true,
				},
			},
		},
		{
			Name:   "list",
			Usage:  "List all plugins",
			Action: plugins.ListPlugins,
		},
	},
}
