package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/urfave/cli/v2"
	"go.codycody31.dev/support/config"
)

func LoadPlugins(app *cli.App) {
	pluginsDir := config.GetConfig().PluginsDir
	err := filepath.Walk(pluginsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".so" {
			pluginName := info.Name()

			// Then trim off _plugin.so
			pluginName = pluginName[:len(pluginName)-10]

			if config.GetConfig().Plugins[pluginName] {
				p, err := plugin.Open(path)
				if err != nil {
					return fmt.Errorf("failed to load plugin %s: %v", pluginName, err)
				}

				symbol, err := p.Lookup("SetupCommands")
				if err != nil {
					return fmt.Errorf("plugin %s does not implement SetupCommands: %v", pluginName, err)
				}

				setupCommands, ok := symbol.(func() []*cli.Command)
				if !ok {
					return fmt.Errorf("invalid SetupCommands signature in plugin %s", pluginName)
				}

				commands := setupCommands()
				app.Commands = append(app.Commands, commands...)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error loading plugins:", err)
	}
}

func RegisterPluginDir(c *cli.Context) error {
	pluginsDir := c.String("dir")

	// Check if the directory exists
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		// Create the directory and any missing parent directories
		err := os.MkdirAll(pluginsDir, 0755)
		if err != nil {
			fmt.Println("Error creating plugin directory:", err)
			return err
		} else {
			fmt.Printf("Created plugin directory %s\n", pluginsDir)
		}
	}

	config.GetConfig().PluginsDir = pluginsDir
	config.SaveConfig()
	fmt.Printf("Plugin directory set to %s\n", pluginsDir)
	return nil
}

func EnablePlugin(c *cli.Context) error {
	pluginName := c.String("name")
	config.GetConfig().Plugins[pluginName] = true
	config.SaveConfig()
	fmt.Printf("Plugin %s enabled\n", pluginName)
	return nil
}

func DisablePlugin(c *cli.Context) error {
	pluginName := c.String("name")
	config.GetConfig().Plugins[pluginName] = false
	config.SaveConfig()
	fmt.Printf("Plugin %s disabled\n", pluginName)
	return nil
}

func ListPlugins(c *cli.Context) error {
	fmt.Println("Plugins:")
	for name, enabled := range config.GetConfig().Plugins {
		fmt.Printf("  %s: %v\n", name, enabled)
	}
	return nil
}
