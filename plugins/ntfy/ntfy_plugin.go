package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
)

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
			},
		},
	}
}

func NtfySend(c *cli.Context) error {
	topic := c.String("topic")
	message := c.String("message")

	if topic == "" || message == "" {
		return fmt.Errorf("both topic and message are required")
	}

	url := fmt.Sprintf("https://ntfy.sh/%s", topic)
	req, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Title", "Notification")

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
