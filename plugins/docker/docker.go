package main

import (
	"fmt"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func Name() string {
	return "docker"
}

func SetupCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "docker",
			Usage: "Manage Docker containers",
			Subcommands: []*cli.Command{
				{
					Name:   "list",
					Usage:  "List all running Docker containers",
					Action: DockerList,
				},
				{
					Name:   "stop",
					Usage:  "Stop a Docker container",
					Action: DockerStop,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "container",
							Aliases:  []string{"c"},
							Usage:    "Container ID or name",
							Required: true,
						},
					},
				},
			},
		},
	}
}

func DockerList(c *cli.Context) error {
	cmd := exec.Command("docker", "ps")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to list containers: %v", err)
	}
	fmt.Println(string(output))
	return nil
}

func DockerStop(c *cli.Context) error {
	container := c.String("container")
	cmd := exec.Command("docker", "stop", container)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop container %s: %v", container, err)
	}
	fmt.Println(string(output))
	return nil
}
