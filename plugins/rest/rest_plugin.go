package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
)

func SetupCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "rest",
			Usage: "Interact with a REST API",
			Subcommands: []*cli.Command{
				{
					Name:   "get",
					Usage:  "Make a GET request",
					Action: RestGet,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "url",
							Aliases:  []string{"u"},
							Usage:    "URL to send the GET request to",
							Required: true,
						},
					},
				},
				{
					Name:   "post",
					Usage:  "Make a POST request",
					Action: RestPost,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "url",
							Aliases:  []string{"u"},
							Usage:    "URL to send the POST request to",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "data",
							Aliases:  []string{"d"},
							Usage:    "Data to send in the POST request",
							Required: true,
						},
					},
				},
			},
		},
	}
}

func RestGet(c *cli.Context) error {
	url := c.String("url")
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Println(string(body))
	return nil
}

func RestPost(c *cli.Context) error {
	url := c.String("url")
	data := c.String("data")
	resp, err := http.Post(url, "application/json", strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Println(string(body))
	return nil
}
