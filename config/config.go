package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Plugins map[string]bool `yaml:"plugins"`
	// PluginsDir     string                            `yaml:"plugins_dir"`
	PluginDirs     []string                          `yaml:"plugin_dirs"`
	PluginSettings map[string]map[string]interface{} `yaml:"plugin_settings"`
}

var configFilePath string
var config Config

func init() {
	configFilePath = getConfigFilePath()
	loadConfig()
}

func getSupportDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		homeDir = "/tmp" // Fallback to /tmp if home directory is not found
	}

	return filepath.Join(homeDir, ".support")
}

func getConfigFilePath() string {
	return filepath.Join(getSupportDir(), "config.yaml")
}

func loadConfig() {
	file, err := os.Open(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			config = Config{
				Plugins:        make(map[string]bool),
				PluginSettings: make(map[string]map[string]interface{}),
				PluginDirs:     []string{filepath.Join(getSupportDir(), "plugins")},
			}
			saveConfig()
			return
		}
		fmt.Println("Error loading config:", err)
		return
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
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

	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(&config)
	if err != nil {
		fmt.Println("Error encoding config:", err)
	}
}

func GetConfig() *Config {
	return &config
}

func SaveConfig() {
	saveConfig()
}

func UpdatePluginSetting(pluginName, key string, value interface{}) error {
	if config.PluginSettings == nil {
		config.PluginSettings = make(map[string]map[string]interface{})
	}

	if _, exists := config.PluginSettings[pluginName]; !exists {
		config.PluginSettings[pluginName] = make(map[string]interface{})
	}

	config.PluginSettings[pluginName][key] = value
	SaveConfig()
	return nil
}

func GetPluginSetting(pluginName, key string) (interface{}, bool) {
	if config.PluginSettings == nil {
		return nil, false
	}

	if pluginSettings, exists := config.PluginSettings[pluginName]; exists {
		value, exists := pluginSettings[key]
		return value, exists
	}

	return nil, false
}
