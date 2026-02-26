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

// contactsCmd is the parent command for contact management.
var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "Manage contacts",
}

// contactsListCmd lists all contacts.
var contactsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all contacts",
	Run:   runContactsList,
}

// contactsGetCmd retrieves a single contact by ID.
var contactsGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a single contact",
	Args:  cobra.ExactArgs(1),
	Run:   runContactsGet,
}

// contactsCreateCmd creates a new contact.
var contactsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new contact",
	Run:   runContactsCreate,
}

// contactsUpdateCmd updates an existing contact.
var contactsUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a contact",
	Args:  cobra.ExactArgs(1),
	Run:   runContactsUpdate,
}

// contactsDeleteCmd deletes a contact.
var contactsDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a contact",
	Args:  cobra.ExactArgs(1),
	Run:   runContactsDelete,
}

// init registers the contacts commands and their flags.
func init() {
	// Create flags.
	contactsCreateCmd.Flags().String("name", "", "Contact name (required)")
	contactsCreateCmd.Flags().String("first-name", "", "First name")
	contactsCreateCmd.Flags().String("last-name", "", "Last name")
	contactsCreateCmd.Flags().String("email", "", "Email address")
	contactsCreateCmd.Flags().String("phone", "", "Phone number")
	contactsCreateCmd.Flags().String("address", "", "Street address")
	contactsCreateCmd.Flags().String("city", "", "City")
	contactsCreateCmd.Flags().String("state", "", "State")
	contactsCreateCmd.Flags().String("zip", "", "Zip code")
	contactsCreateCmd.Flags().String("country", "", "Country")
	contactsCreateCmd.Flags().String("website", "", "Website URL")
	contactsCreateCmd.MarkFlagRequired("name")

	// Update flags.
	contactsUpdateCmd.Flags().String("name", "", "Contact name")
	contactsUpdateCmd.Flags().String("first-name", "", "First name")
	contactsUpdateCmd.Flags().String("last-name", "", "Last name")
	contactsUpdateCmd.Flags().String("email", "", "Email address")
	contactsUpdateCmd.Flags().String("phone", "", "Phone number")
	contactsUpdateCmd.Flags().String("address", "", "Street address")
	contactsUpdateCmd.Flags().String("city", "", "City")
	contactsUpdateCmd.Flags().String("state", "", "State")
	contactsUpdateCmd.Flags().String("zip", "", "Zip code")
	contactsUpdateCmd.Flags().String("country", "", "Country")
	contactsUpdateCmd.Flags().String("website", "", "Website URL")

	// List flags.
	contactsListCmd.Flags().String("search", "", "Search contacts by name")

	contactsCmd.AddCommand(contactsListCmd)
	contactsCmd.AddCommand(contactsGetCmd)
	contactsCmd.AddCommand(contactsCreateCmd)
	contactsCmd.AddCommand(contactsUpdateCmd)
	contactsCmd.AddCommand(contactsDeleteCmd)
	rootCmd.AddCommand(contactsCmd)
}

// runContactsList fetches and displays all contacts.
func runContactsList(cmd *cobra.Command, args []string) {
	client := newClient()

	params := map[string]string{}
	search, _ := cmd.Flags().GetString("search")
	if search != "" {
		params["search"] = search
	}

	contacts, err := client.GetContacts(params)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(contacts)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tEMAIL\tPHONE")
	for _, c := range contacts {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", c.ID, c.Name, c.Email, c.Phone)
	}
	w.Flush()
}

// runContactsGet fetches and displays a single contact.
func runContactsGet(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid contact ID")
		os.Exit(1)
	}

	contact, err := client.GetContact(uint(id))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(contact)
		return
	}

	fmt.Printf("ID:      %d\n", contact.ID)
	fmt.Printf("Name:    %s\n", contact.Name)
	if contact.FirstName != "" {
		fmt.Printf("First:   %s\n", contact.FirstName)
	}
	if contact.LastName != "" {
		fmt.Printf("Last:    %s\n", contact.LastName)
	}
	if contact.Email != "" {
		fmt.Printf("Email:   %s\n", contact.Email)
	}
	if contact.Phone != "" {
		fmt.Printf("Phone:   %s\n", contact.Phone)
	}
	if contact.Address != "" {
		fmt.Printf("Address: %s\n", contact.Address)
	}
	if contact.City != "" {
		fmt.Printf("City:    %s\n", contact.City)
	}
	if contact.State != "" {
		fmt.Printf("State:   %s\n", contact.State)
	}
	if contact.Zip != "" {
		fmt.Printf("Zip:     %s\n", contact.Zip)
	}
	if contact.Country != "" {
		fmt.Printf("Country: %s\n", contact.Country)
	}
	if contact.Website != "" {
		fmt.Printf("Website: %s\n", contact.Website)
	}
}

// runContactsCreate creates a new contact from flags.
func runContactsCreate(cmd *cobra.Command, args []string) {
	client := newClient()

	name, _ := cmd.Flags().GetString("name")
	firstName, _ := cmd.Flags().GetString("first-name")
	lastName, _ := cmd.Flags().GetString("last-name")
	email, _ := cmd.Flags().GetString("email")
	phone, _ := cmd.Flags().GetString("phone")
	address, _ := cmd.Flags().GetString("address")
	city, _ := cmd.Flags().GetString("city")
	state, _ := cmd.Flags().GetString("state")
	zip, _ := cmd.Flags().GetString("zip")
	country, _ := cmd.Flags().GetString("country")
	website, _ := cmd.Flags().GetString("website")

	contact, err := client.CreateContact(&api.ContactCreateRequest{
		Name:      name,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Address:   address,
		City:      city,
		State:     state,
		Zip:       zip,
		Country:   country,
		Website:   website,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(contact)
		return
	}

	fmt.Printf("Created contact %d: %s\n", contact.ID, contact.Name)
}

// runContactsUpdate updates an existing contact from flags.
func runContactsUpdate(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid contact ID")
		os.Exit(1)
	}

	req := &api.ContactUpdateRequest{}

	if cmd.Flags().Changed("name") {
		v, _ := cmd.Flags().GetString("name")
		req.Name = v
	}
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
	if cmd.Flags().Changed("phone") {
		v, _ := cmd.Flags().GetString("phone")
		req.Phone = v
	}
	if cmd.Flags().Changed("address") {
		v, _ := cmd.Flags().GetString("address")
		req.Address = v
	}
	if cmd.Flags().Changed("city") {
		v, _ := cmd.Flags().GetString("city")
		req.City = v
	}
	if cmd.Flags().Changed("state") {
		v, _ := cmd.Flags().GetString("state")
		req.State = v
	}
	if cmd.Flags().Changed("zip") {
		v, _ := cmd.Flags().GetString("zip")
		req.Zip = v
	}
	if cmd.Flags().Changed("country") {
		v, _ := cmd.Flags().GetString("country")
		req.Country = v
	}
	if cmd.Flags().Changed("website") {
		v, _ := cmd.Flags().GetString("website")
		req.Website = v
	}

	contact, err := client.UpdateContact(uint(id), req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(contact)
		return
	}

	fmt.Printf("Updated contact %d: %s\n", contact.ID, contact.Name)
}

// runContactsDelete deletes a contact by ID.
func runContactsDelete(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid contact ID")
		os.Exit(1)
	}

	if err := client.DeleteContact(uint(id)); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted contact %d\n", id)
}
