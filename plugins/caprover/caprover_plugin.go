package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
	"go.codycody31.dev/support/config"
)

func Name() string {
	return "caprover"
}

func SetupCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "caprover",
			Usage: "Manage CapRover instances",
			Subcommands: []*cli.Command{
				{
					Name:   "configure",
					Usage:  "Configure CapRover",
					Action: CaproverConfigure,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "url",
							Aliases:  []string{"u"},
							Usage:    "CapRover server URL",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "password",
							Aliases:  []string{"p"},
							Usage:    "CapRover password",
							Required: true,
						},
					},
				},
				// Bulk delete apps via regex
				{
					Name:   "delete",
					Usage:  "Delete apps",
					Action: CaproverDelete,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "regex",
							Aliases:  []string{"r"},
							Usage:    "Regex to match app names",
							Required: true,
						},
						&cli.BoolFlag{
							Name:        "dry-run",
							Usage:       "Dry run",
							DefaultText: "false",
						},
					},
				},
			},
		},
	}
}

func CaproverConfigure(c *cli.Context) error {
	url := c.String("url")
	password := c.String("password")

	// Strip the trailing slash from the URL
	url = strings.TrimSuffix(url, "/")

	// Body for the request
	body := strings.NewReader(`{"password": "` + password + `"}`)

	// Create the request
	req, err := http.NewRequest("POST", url+"/api/v2/login", body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login: %v", resp.Status)
	}

	// Read the response body
	response := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	token := response["data"].(map[string]interface{})["token"].(string)

	// Store the CapRover server URL and token
	err = config.UpdatePluginSetting("caprover", "server", url)
	if err != nil {
		return fmt.Errorf("failed to set caprover server: %v", err)
	}
	err = config.UpdatePluginSetting("caprover", "token", token)
	if err != nil {
		return fmt.Errorf("failed to set caprover token: %v", err)
	}

	fmt.Printf("CapRover server set to %s\n", url)

	return nil
}

func CaproverDelete(c *cli.Context) error {
	regex := c.String("regex")
	dryRun := c.Bool("dry-run")
	matched := 0

	// Get the CapRover server URL and token
	server, exists := config.GetPluginSetting("caprover", "server")
	if !exists {
		return fmt.Errorf("caprover server not set")
	}
	token, exists := config.GetPluginSetting("caprover", "token")
	if !exists {
		return fmt.Errorf("caprover token not set")
	}

	// Create the request
	req, err := http.NewRequest("GET", server.(string)+"/api/v2/user/apps/appDefinitions", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-captain-auth", token.(string))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get apps: %v", resp.Status)
	}

	// Read the response body
	response := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	if response["status"].(float64) == 1106 {
		return fmt.Errorf(response["description"].(string))
	}

	// Loop through the apps
	apps := response["data"].(map[string]interface{})["appDefinitions"].([]interface{})
	totalApps := len(apps)

	fmt.Printf("Found %d apps\n", totalApps)

	for _, app := range apps {
		appName := app.(map[string]interface{})["appName"].(string)

		// Check if the app matches the regex
		// TODO: Allow setting the regex flags
		// TODO: Along with strict matching like ^app$, etc being all configurable via flags
		if !strings.Contains(appName, regex) {
			continue
		}

		matched++

		if dryRun {
			fmt.Printf("Would delete app: %s\n", appName)
			continue
		}

		// Create the request
		req, err := http.NewRequest("POST", server.(string)+"/api/v2/user/apps/appDefinitions/delete", strings.NewReader(`{"appName": "`+appName+`"}`))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-captain-auth", token.(string))

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to delete app: %v", resp.Status)
		}

		fmt.Printf("Deleted app: %s\n", appName)
	}

	if matched == 0 {
		fmt.Println("No apps matched the regex")
	} else {
		fmt.Printf("Deleted %d apps\n", matched)
	}

	return nil
}
