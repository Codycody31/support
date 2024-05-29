package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"runtime"

	"github.com/urfave/cli/v2"
)

var configFilePath string
var enabledPlugins = map[string]bool{}

func init() {
	configFilePath = getConfigFilePath()
	loadConfig()
}

func getConfigFilePath() string {
	var configDir string

	if runtime.GOOS == "windows" {
		configDir = os.Getenv("APPDATA")
	} else {
		configDir = "/etc"
	}

	return filepath.Join(configDir, "support", "plugins.json")
}

func loadConfig() {
	file, err := os.Open(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create the folder if it doesn't exist
			err = os.MkdirAll(filepath.Dir(configFilePath), 0755)
			if err != nil {
				fmt.Println("Error creating config folder:", err)
				return
			}

			saveConfig()
			return
		}
		fmt.Println("Error loading config:", err)
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&enabledPlugins)
	if err != nil {
		fmt.Println("Error decoding config:", err)
	}
}

func saveConfig() {
	file, err := os.Create(configFilePath)
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(&enabledPlugins)
	if err != nil {
		fmt.Println("Error encoding config:", err)
	}
}

func LoadPlugins(app *cli.App) {
	pluginsDir := "./plugins_dir"
	err := filepath.Walk(pluginsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".so" {
			pluginName := info.Name()

			// Then trim off _plugin.so
			pluginName = pluginName[:len(pluginName)-10]

			if enabledPlugins[pluginName] {
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

func EnablePlugin(c *cli.Context) error {
	pluginName := c.String("name")
	enabledPlugins[pluginName] = true
	saveConfig()
	fmt.Printf("Plugin %s enabled\n", pluginName)
	return nil
}

func DisablePlugin(c *cli.Context) error {
	pluginName := c.String("name")
	enabledPlugins[pluginName] = false
	saveConfig()
	fmt.Printf("Plugin %s disabled\n", pluginName)
	return nil
}

func ListPlugins(c *cli.Context) error {
	fmt.Println("Plugins:")
	for name, enabled := range enabledPlugins {
		fmt.Printf("  %s: %v\n", name, enabled)
	}
	return nil
}
