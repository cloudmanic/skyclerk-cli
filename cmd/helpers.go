// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloudmanic/skyclerk-cli/internal/api"
	"github.com/cloudmanic/skyclerk-cli/internal/config"
)

// newClient loads the config and creates a new authenticated API client.
func newClient() *api.Client {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	accountID := cfg.DefaultAccountID
	if accountOverride > 0 {
		accountID = accountOverride
	}

	if accountID == 0 {
		fmt.Fprintln(os.Stderr, "Error: no account selected. Run 'skyclerk accounts use <id>' first.")
		os.Exit(1)
	}

	baseURL := cfg.ApiURL
	if baseURL == "" {
		baseURL = config.DefaultApiURL
	}

	return api.NewClient(baseURL, cfg.AccessToken, accountID)
}

// newClientNoAccount loads the config and creates a client without requiring an account ID.
func newClientNoAccount() (*api.Client, *config.Config) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	baseURL := cfg.ApiURL
	if baseURL == "" {
		baseURL = config.DefaultApiURL
	}

	return api.NewClient(baseURL, cfg.AccessToken, 0), cfg
}

// printJSON pretty-prints any value as indented JSON.
func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error formatting JSON:", err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}
