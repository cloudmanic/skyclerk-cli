// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package cmd

import (
	"fmt"
	"os"

	"github.com/cloudmanic/skyclerk-cli/internal/api"
	"github.com/spf13/cobra"
)

// meCmd displays the current user profile.
var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Show your user profile",
	Run:   runMe,
}

// meUpdateCmd updates the current user profile.
var meUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update your user profile",
	Run:   runMeUpdate,
}

// init registers the me commands and their flags.
func init() {
	meUpdateCmd.Flags().String("first-name", "", "First name")
	meUpdateCmd.Flags().String("last-name", "", "Last name")
	meUpdateCmd.Flags().String("email", "", "Email address")

	meCmd.AddCommand(meUpdateCmd)
	rootCmd.AddCommand(meCmd)
}

// runMe fetches and displays the current user profile.
func runMe(cmd *cobra.Command, args []string) {
	client := newClient()

	me, err := client.GetMe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(me)
		return
	}

	fmt.Printf("ID:         %d\n", me.ID)
	fmt.Printf("First Name: %s\n", me.FirstName)
	fmt.Printf("Last Name:  %s\n", me.LastName)
	fmt.Printf("Email:      %s\n", me.Email)
	fmt.Printf("Status:     %s\n", me.Status)
}

// runMeUpdate updates the current user profile from flags.
func runMeUpdate(cmd *cobra.Command, args []string) {
	client := newClient()

	req := &api.MeUpdateRequest{}

	if cmd.Flags().Changed("first-name") {
		v, _ := cmd.Flags().GetString("first-name")
		req.FirstName = v
	}
	if cmd.Flags().Changed("last-name") {
		v, _ := cmd.Flags().GetString("last-name")
		req.LastName = v
	}
	if cmd.Flags().Changed("email") {
		v, _ := cmd.Flags().GetString("email")
		req.Email = v
	}

	me, err := client.UpdateMe(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(me)
		return
	}

	fmt.Printf("Updated profile: %s %s (%s)\n", me.FirstName, me.LastName, me.Email)
}
