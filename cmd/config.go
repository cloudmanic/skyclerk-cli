// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package cmd

import (
	"fmt"
	"os"

	"github.com/cloudmanic/skyclerk-cli/internal/config"
	"github.com/spf13/cobra"
)

// configCmd is the parent command for config management.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

// configShowCmd displays the current configuration with masked secrets.
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run:   runConfigShow,
}

// configInitCmd allows manual configuration setup.
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Manually initialize configuration",
	Run:   runConfigInit,
}

// init registers the config commands.
func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)
	rootCmd.AddCommand(configCmd)
}

// runConfigShow displays the current configuration.
func runConfigShow(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(map[string]interface{}{
			"access_token":       config.MaskString(cfg.AccessToken),
			"user_id":            cfg.UserID,
			"default_account_id": cfg.DefaultAccountID,
			"api_url":            cfg.ApiURL,
		})
		return
	}

	fmt.Printf("Access Token:       %s\n", config.MaskString(cfg.AccessToken))
	fmt.Printf("User ID:            %d\n", cfg.UserID)
	fmt.Printf("Default Account ID: %d\n", cfg.DefaultAccountID)
	fmt.Printf("API URL:            %s\n", cfg.ApiURL)
}

// runConfigInit prompts the user to manually set config values.
func runConfigInit(cmd *cobra.Command, args []string) {
	var accessToken, apiURL string
	var accountID uint

	fmt.Print("Access Token: ")
	fmt.Scanln(&accessToken)

	fmt.Printf("API URL (default: %s): ", config.DefaultApiURL)
	fmt.Scanln(&apiURL)
	if apiURL == "" {
		apiURL = config.DefaultApiURL
	}

	fmt.Print("Default Account ID: ")
	fmt.Scanln(&accountID)

	cfg := &config.Config{
		AccessToken:      accessToken,
		DefaultAccountID: accountID,
		ApiURL:           apiURL,
	}

	if err := config.Save(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving config:", err)
		os.Exit(1)
	}

	fmt.Println("Configuration saved.")
}
