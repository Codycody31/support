package main

import (
	"os"

	"go.codycody31.dev/support/plugins"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "support",
		Usage: "System Utilities and Plugin-based Operations, Routines, and Tasks",
		Commands: []*cli.Command{
			PluginsCommand,
		},
	}

	plugins.LoadPlugins(app)

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
