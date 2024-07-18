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
	for _, pluginsDir := range config.GetConfig().PluginDirs {
		err := filepath.Walk(pluginsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(path) == ".so" {
				pluginName := info.Name()

				// Then trim off .so
				pluginName = pluginName[:len(pluginName)-3]

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
			fmt.Printf("Error loading plugins from directory %s: %v\n", pluginsDir, err)
		}
	}
}

func RegisterPluginDir(c *cli.Context) error {
	pluginsDir := c.Args().First()

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

	// Check if the directory is already registered
	for _, dir := range config.GetConfig().PluginDirs {
		if dir == pluginsDir {
			fmt.Printf("Plugin directory %s is already registered\n", pluginsDir)
			return nil
		}
	}

	config.GetConfig().PluginDirs = append(config.GetConfig().PluginDirs, pluginsDir)
	config.SaveConfig()
	fmt.Printf("Plugin directory set to %s\n", pluginsDir)
	return nil
}

func EnablePlugin(c *cli.Context) error {
	pluginName := c.Args().First()

	configData := config.GetConfig()

	// Initialize the Plugins map if it's nil
	if configData.Plugins == nil {
		configData.Plugins = make(map[string]bool)
	}

	if _, exists := configData.Plugins[pluginName]; !exists {
		pluginExists := false

		// Verify the plugin exists (by checking for the .so file)
		for _, pluginsDir := range configData.PluginDirs {
			pluginPath := filepath.Join(pluginsDir, pluginName+".so")
			if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
				pluginExists = true
			}
		}

		if !pluginExists {
			fmt.Printf("Plugin %s does not exist\n", pluginName)
			return nil
		}

		// Add the plugin to the config
		configData.Plugins[pluginName] = true
	}

	config.SaveConfig()
	fmt.Printf("Plugin %s enabled\n", pluginName)
	return nil
}

func DisablePlugin(c *cli.Context) error {
	pluginName := c.Args().First()
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
