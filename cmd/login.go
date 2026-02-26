// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/cloudmanic/skyclerk-cli/internal/api"
	"github.com/cloudmanic/skyclerk-cli/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// loginCmd authenticates a user with email and password and stores the access token.
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Skyclerk with your email and password",
	Run:   runLogin,
}

// logoutCmd revokes the access token and removes the config file.
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out and revoke your access token",
	Run:   runLogout,
}

// init registers the login and logout commands.
func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
}

// runLogin handles the login command execution.
func runLogin(cmd *cobra.Command, args []string) {
	apiURL := config.DefaultApiURL

	// Prompt for client ID.
	fmt.Print("Client ID: ")
	var clientID string
	fmt.Scanln(&clientID)

	if clientID == "" {
		fmt.Fprintln(os.Stderr, "Error: client ID is required")
		os.Exit(1)
	}

	// Prompt for email.
	fmt.Print("Email: ")
	var email string
	fmt.Scanln(&email)

	if email == "" {
		fmt.Fprintln(os.Stderr, "Error: email is required")
		os.Exit(1)
	}

	// Prompt for password (hidden input).
	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading password:", err)
		os.Exit(1)
	}

	password := string(passwordBytes)
	if password == "" {
		fmt.Fprintln(os.Stderr, "Error: password is required")
		os.Exit(1)
	}

	// Authenticate with the API.
	client := api.NewClient(apiURL, "", 0)
	resp, err := client.Login(email, password, clientID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	// Store the token temporarily to fetch user accounts.
	client = api.NewClient(apiURL, resp.AccessToken, 0)
	user, err := client.GetAuthUser()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error fetching user profile:", err)
		os.Exit(1)
	}

	// Determine the default account.
	var defaultAccountID uint
	if len(user.Accounts) == 1 {
		defaultAccountID = user.Accounts[0].ID
		fmt.Printf("Using account: %s (ID: %d)\n", user.Accounts[0].Name, defaultAccountID)
	} else if len(user.Accounts) > 1 {
		fmt.Println("\nAvailable accounts:")
		for _, acct := range user.Accounts {
			fmt.Printf("  [%d] %s\n", acct.ID, acct.Name)
		}
		fmt.Print("\nSelect default account ID: ")
		fmt.Scanln(&defaultAccountID)

		// Validate the selected account belongs to this user.
		valid := false
		for _, acct := range user.Accounts {
			if acct.ID == defaultAccountID {
				valid = true
				break
			}
		}
		if !valid {
			fmt.Fprintln(os.Stderr, "Error: invalid account ID")
			os.Exit(1)
		}
	}

	// Save the config.
	cfg := &config.Config{
		AccessToken:      resp.AccessToken,
		UserID:           resp.UserID,
		DefaultAccountID: defaultAccountID,
		ApiURL:           apiURL,
		ClientID:         clientID,
	}

	if err := config.Save(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving config:", err)
		os.Exit(1)
	}

	fmt.Printf("\nLogged in as %s %s (%s)\n", user.FirstName, user.LastName, user.Email)
}

// runLogout handles the logout command execution.
func runLogout(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	baseURL := cfg.ApiURL
	if baseURL == "" {
		baseURL = config.DefaultApiURL
	}

	// Revoke the token on the server.
	client := api.NewClient(baseURL, cfg.AccessToken, 0)
	if err := client.Logout(); err != nil {
		fmt.Fprintln(os.Stderr, "Warning: could not revoke token:", err)
	}

	// Delete the local config file.
	if err := config.Delete(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Println("Logged out successfully.")
}
