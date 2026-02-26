// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/cloudmanic/skyclerk-cli/internal/api"
	"github.com/spf13/cobra"
)

// usersCmd is the parent command for user management.
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage account users and invitations",
}

// usersListCmd lists all users in the current account.
var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users in the account",
	Run:   runUsersList,
}

// usersRemoveCmd removes a user from the current account.
var usersRemoveCmd = &cobra.Command{
	Use:   "remove [id]",
	Short: "Remove a user from the account",
	Args:  cobra.ExactArgs(1),
	Run:   runUsersRemove,
}

// usersInviteCmd sends an invitation to join the account.
var usersInviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Invite a user to the account",
	Run:   runUsersInvite,
}

// usersInvitesCmd lists pending invitations.
var usersInvitesCmd = &cobra.Command{
	Use:   "invites",
	Short: "List pending invitations",
	Run:   runUsersInvites,
}

// usersCancelInviteCmd cancels a pending invitation.
var usersCancelInviteCmd = &cobra.Command{
	Use:   "cancel-invite [id]",
	Short: "Cancel a pending invitation",
	Args:  cobra.ExactArgs(1),
	Run:   runUsersCancelInvite,
}

// init registers the users commands and their flags.
func init() {
	// Invite flags.
	usersInviteCmd.Flags().String("email", "", "Email address to invite")
	usersInviteCmd.Flags().String("first-name", "", "First name")
	usersInviteCmd.Flags().String("last-name", "", "Last name")
	usersInviteCmd.Flags().String("message", "", "Invitation message")
	usersInviteCmd.MarkFlagRequired("email")
	usersInviteCmd.MarkFlagRequired("first-name")
	usersInviteCmd.MarkFlagRequired("last-name")

	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersRemoveCmd)
	usersCmd.AddCommand(usersInviteCmd)
	usersCmd.AddCommand(usersInvitesCmd)
	usersCmd.AddCommand(usersCancelInviteCmd)
	rootCmd.AddCommand(usersCmd)
}

// runUsersList fetches and displays all users in the account.
func runUsersList(cmd *cobra.Command, args []string) {
	client := newClient()

	users, err := client.GetUsers()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(users)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tEMAIL\tSTATUS")
	for _, u := range users {
		fmt.Fprintf(w, "%d\t%s %s\t%s\t%s\n", u.ID, u.FirstName, u.LastName, u.Email, u.Status)
	}
	w.Flush()
}

// runUsersRemove removes a user from the account.
func runUsersRemove(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid user ID")
		os.Exit(1)
	}

	if err := client.RemoveUser(uint(id)); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Removed user %d from account\n", id)
}

// runUsersInvite sends an invitation to join the account.
func runUsersInvite(cmd *cobra.Command, args []string) {
	client := newClient()

	email, _ := cmd.Flags().GetString("email")
	firstName, _ := cmd.Flags().GetString("first-name")
	lastName, _ := cmd.Flags().GetString("last-name")
	message, _ := cmd.Flags().GetString("message")

	invite, err := client.CreateInvite(&api.InviteCreateRequest{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Message:   message,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(invite)
		return
	}

	fmt.Printf("Invitation sent to %s (expires: %s)\n", invite.Email, invite.ExpiresAt)
}

// runUsersInvites lists pending invitations.
func runUsersInvites(cmd *cobra.Command, args []string) {
	client := newClient()

	invites, err := client.GetInvites()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(invites)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tEMAIL\tNAME\tEXPIRES")
	for _, inv := range invites {
		fmt.Fprintf(w, "%d\t%s\t%s %s\t%s\n", inv.ID, inv.Email, inv.FirstName, inv.LastName, inv.ExpiresAt)
	}
	w.Flush()
}

// runUsersCancelInvite cancels a pending invitation.
func runUsersCancelInvite(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid invite ID")
		os.Exit(1)
	}

	if err := client.CancelInvite(uint(id)); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Cancelled invitation %d\n", id)
}
