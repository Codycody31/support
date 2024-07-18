package main

import (
	"go.codycody31.dev/support/plugins"

	"github.com/urfave/cli/v2"
)

var PluginsCommand = &cli.Command{
	Name:  "plugins",
	Usage: "Manage plugins",
	Subcommands: []*cli.Command{
		{
			Name:   "register",
			Usage:  "Register a plugin directory",
			Action: plugins.RegisterPluginDir,
		},
		{
			Name:   "enable",
			Usage:  "Enable a plugin",
			Action: plugins.EnablePlugin,
		},
		{
			Name:   "disable",
			Usage:  "Disable a plugin",
			Action: plugins.DisablePlugin,
		},
		{
			Name:   "list",
			Usage:  "List all plugins",
			Action: plugins.ListPlugins,
		},
	},
}
