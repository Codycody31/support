package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
	"go.codycody31.dev/support/config"
)

func Name() string {
	return "ntfy"
}

func SetupCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "ntfy",
			Usage: "Send notifications via ntfy.sh",
			Subcommands: []*cli.Command{
				{
					Name:   "send",
					Usage:  "Send a notification",
					Action: NtfySend,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "topic",
							Aliases:  []string{"t"},
							Usage:    "Notification topic",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "message",
							Aliases:  []string{"m"},
							Usage:    "Notification message",
							Required: true,
						},
					},
				},
				{
					Name:   "configure",
					Usage:  "Configure ntfy",
					Action: ConfigureNtfy,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "url",
							Aliases:  []string{"u"},
							Usage:    "Ntfy server URL",
							Required: true,
						},
						&cli.StringFlag{
							Name:    "access-token",
							Aliases: []string{"a"},
							Usage:   "Ntfy access token",
						},
					},
				},
			},
		},
	}
}

func NtfySend(c *cli.Context) error {
	topic := c.String("topic")
	message := c.String("message")
	server, _ := config.GetPluginSetting("ntfy", "server")

	if topic == "" || message == "" {
		return fmt.Errorf("both topic and message are required")
	}

	serverURL := "https://ntfy.sh"
	if server != nil {
		serverURL = server.(string)
	}

	url := fmt.Sprintf("%s/%s", serverURL, topic)
	req, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Title", "Notification")
	if token, _ := config.GetPluginSetting("ntfy", "access-token"); token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	fmt.Println("Notification sent successfully!")
	return nil
}

func ConfigureNtfy(c *cli.Context) error {
	url := c.String("url")
	accessToken := c.String("access-token")

	// Strip the trailing slash from the URL
	url = strings.TrimSuffix(url, "/")

	err := config.UpdatePluginSetting("ntfy", "server", url)
	if err != nil {
		return fmt.Errorf("failed to set ntfy server: %v", err)
	}

	if accessToken != "" {
		err = config.UpdatePluginSetting("ntfy", "access-token", accessToken)
		if err != nil {
			return fmt.Errorf("failed to set ntfy access token: %v", err)
		}
	}

	fmt.Printf("Ntfy server set to %s\n", url)
	return nil
}
