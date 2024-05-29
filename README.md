# SUPPORT

**SUPPORT** (System Utilities and Plugin-based Operations, Routines, and Tasks) is a CLI utility tool written in Go. It is designed to provide a flexible and extendable way to perform various system tasks through plugins. The tool uses the `urfave/cli/v2` package to manage CLI commands and dynamically loads plugins to extend its functionality.

## Features

- **Plugin-Based Architecture**: Easily extend functionality by adding or removing plugins.
- **Dynamic Plugin Loading**: Load plugins dynamically from a specified directory.
- **Plugin Management**: Enable or disable plugins using CLI commands, with configuration stored in a JSON file.

## Getting Started

### Prerequisites

- Go 1.22.2 or later

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/Codycody31/support.git
    cd support
    ```

2. Build all plugins using the provided script:

    ```sh
    ./build_plugins.sh
    ```

3. Build the core application:

    ```sh
    go build -o support main.go
    ```

### Usage

1. Run the `support` CLI tool:

    ```sh
    ./support
    ```

2. Use the `plugin` command to manage plugins:
    - Enable a plugin:

      ```sh
      ./support plugin enable --name ntfy_plugin.so
      ```

    - Disable a plugin:

      ```sh
      ./support plugin disable --name ntfy_plugin.so
      ```

3. Use the dynamically loaded plugin commands:
    - For example, to send a notification using the `ntfy` plugin:

      ```sh
      ./support ntfy send --topic 'mytopic' --message 'Hello, World!'
      ```

## Plugin Development

### Creating a Plugin

1. Create a new directory for your plugin inside the `plugins` directory:

    ```sh
    mkdir plugins/my_plugin
    ```

2. Implement your plugin, ensuring it has a `SetupCommands` function that returns the commands to be added to the CLI:

    ```go
    // plugins/my_plugin/my_plugin.go
    package main

    import (
        "github.com/urfave/cli/v2"
    )

    func SetupCommands() []*cli.Command {
        return []*cli.Command{
            {
                Name:  "mycommand",
                Usage: "My custom command",
                Action: func(c *cli.Context) error {
                    // Your command implementation here
                    return nil
                },
            },
        }
    }
    ```

3. Build all plugins using the provided script:

    ```sh
    ./build_plugins.sh
    ```

4. Enable the plugin and use the new commands as shown in the usage section.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss your ideas or improvements.

## License

This project is licensed under the MIT License.

## Acknowledgements

- [urfave/cli](https://github.com/urfave/cli) for the CLI framework.
