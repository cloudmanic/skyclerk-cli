// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/cloudmanic/skyclerk-cli/internal/config"
	"github.com/spf13/cobra"
)

// accountsCmd is the parent command for account management.
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage Skyclerk accounts",
}

// accountsListCmd lists all accounts the user belongs to.
var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts you belong to",
	Run:   runAccountsList,
}

// accountsShowCmd shows the current account details.
var accountsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current account details",
	Run:   runAccountsShow,
}

// accountsUseCmd sets the default account ID.
var accountsUseCmd = &cobra.Command{
	Use:   "use [id]",
	Short: "Set the default account ID",
	Args:  cobra.ExactArgs(1),
	Run:   runAccountsUse,
}

// init registers the accounts commands.
func init() {
	accountsCmd.AddCommand(accountsListCmd)
	accountsCmd.AddCommand(accountsShowCmd)
	accountsCmd.AddCommand(accountsUseCmd)
	rootCmd.AddCommand(accountsCmd)
}

// runAccountsList fetches and displays all user accounts.
func runAccountsList(cmd *cobra.Command, args []string) {
	client, _ := newClientNoAccount()

	user, err := client.GetAuthUser()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(user.Accounts)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCURRENCY\tLOCALE")
	for _, acct := range user.Accounts {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", acct.ID, acct.Name, acct.Currency, acct.Locale)
	}
	w.Flush()
}

// runAccountsShow fetches and displays the current account.
func runAccountsShow(cmd *cobra.Command, args []string) {
	client := newClient()

	account, err := client.GetAccount()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(account)
		return
	}

	fmt.Printf("ID:       %d\n", account.ID)
	fmt.Printf("Name:     %s\n", account.Name)
	fmt.Printf("Currency: %s\n", account.Currency)
	fmt.Printf("Locale:   %s\n", account.Locale)

	if account.Address != "" {
		fmt.Printf("Address:  %s\n", account.Address)
	}
	if account.City != "" {
		fmt.Printf("City:     %s\n", account.City)
	}
	if account.State != "" {
		fmt.Printf("State:    %s\n", account.State)
	}
	if account.Zip != "" {
		fmt.Printf("Zip:      %s\n", account.Zip)
	}
	if account.Country != "" {
		fmt.Printf("Country:  %s\n", account.Country)
	}
}

// runAccountsUse sets the default account ID in the config.
func runAccountsUse(cmd *cobra.Command, args []string) {
	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid account ID")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	cfg.DefaultAccountID = uint(id)

	if err := config.Save(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving config:", err)
		os.Exit(1)
	}

	fmt.Printf("Default account set to %d\n", id)
}
