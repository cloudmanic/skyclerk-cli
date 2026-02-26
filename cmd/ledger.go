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

// ledgerCmd is the parent command for ledger management.
var ledgerCmd = &cobra.Command{
	Use:   "ledger",
	Short: "Manage ledger entries",
}

// ledgerListCmd lists ledger entries with pagination.
var ledgerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ledger entries",
	Run:   runLedgerList,
}

// ledgerGetCmd retrieves a single ledger entry by ID.
var ledgerGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a single ledger entry",
	Args:  cobra.ExactArgs(1),
	Run:   runLedgerGet,
}

// ledgerCreateCmd creates a new ledger entry.
var ledgerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new ledger entry",
	Run:   runLedgerCreate,
}

// ledgerUpdateCmd updates an existing ledger entry.
var ledgerUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a ledger entry",
	Args:  cobra.ExactArgs(1),
	Run:   runLedgerUpdate,
}

// ledgerDeleteCmd deletes a ledger entry.
var ledgerDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a ledger entry",
	Args:  cobra.ExactArgs(1),
	Run:   runLedgerDelete,
}

// ledgerSummaryCmd displays a ledger summary.
var ledgerSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show ledger summary (income, expense, profit)",
	Run:   runLedgerSummary,
}

// init registers the ledger commands and their flags.
func init() {
	// List flags.
	ledgerListCmd.Flags().String("limit", "25", "Number of entries to return")
	ledgerListCmd.Flags().String("page", "1", "Page number")
	ledgerListCmd.Flags().String("sort", "DESC", "Sort direction (ASC or DESC)")

	// Create flags.
	ledgerCreateCmd.Flags().Float64("amount", 0, "Transaction amount (negative for expense)")
	ledgerCreateCmd.Flags().String("date", "", "Transaction date (YYYY-MM-DD)")
	ledgerCreateCmd.Flags().Uint("contact-id", 0, "Contact ID")
	ledgerCreateCmd.Flags().Uint("category-id", 0, "Category ID")
	ledgerCreateCmd.Flags().String("note", "", "Transaction note")
	ledgerCreateCmd.MarkFlagRequired("amount")
	ledgerCreateCmd.MarkFlagRequired("date")
	ledgerCreateCmd.MarkFlagRequired("contact-id")
	ledgerCreateCmd.MarkFlagRequired("category-id")

	// Update flags.
	ledgerUpdateCmd.Flags().Float64("amount", 0, "Transaction amount")
	ledgerUpdateCmd.Flags().String("date", "", "Transaction date (YYYY-MM-DD)")
	ledgerUpdateCmd.Flags().Uint("contact-id", 0, "Contact ID")
	ledgerUpdateCmd.Flags().Uint("category-id", 0, "Category ID")
	ledgerUpdateCmd.Flags().String("note", "", "Transaction note")

	ledgerCmd.AddCommand(ledgerListCmd)
	ledgerCmd.AddCommand(ledgerGetCmd)
	ledgerCmd.AddCommand(ledgerCreateCmd)
	ledgerCmd.AddCommand(ledgerUpdateCmd)
	ledgerCmd.AddCommand(ledgerDeleteCmd)
	ledgerCmd.AddCommand(ledgerSummaryCmd)
	rootCmd.AddCommand(ledgerCmd)
}

// runLedgerList fetches and displays ledger entries.
func runLedgerList(cmd *cobra.Command, args []string) {
	client := newClient()

	limit, _ := cmd.Flags().GetString("limit")
	page, _ := cmd.Flags().GetString("page")
	sort, _ := cmd.Flags().GetString("sort")

	params := map[string]string{
		"limit": limit,
		"page":  page,
		"sort":  sort,
	}

	ledgers, err := client.GetLedgers(params)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(ledgers)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tDATE\tAMOUNT\tCONTACT\tCATEGORY\tNOTE")
	for _, l := range ledgers {
		fmt.Fprintf(w, "%d\t%s\t%.2f\t%s\t%s\t%s\n",
			l.ID, l.Date, l.Amount, l.Contact.Name, l.Category.Name, l.Note)
	}
	w.Flush()
}

// runLedgerGet fetches and displays a single ledger entry.
func runLedgerGet(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid ledger ID")
		os.Exit(1)
	}

	ledger, err := client.GetLedger(uint(id))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(ledger)
		return
	}

	fmt.Printf("ID:        %d\n", ledger.ID)
	fmt.Printf("Date:      %s\n", ledger.Date)
	fmt.Printf("Amount:    %.2f\n", ledger.Amount)
	fmt.Printf("Contact:   %s\n", ledger.Contact.Name)
	fmt.Printf("Category:  %s\n", ledger.Category.Name)
	if ledger.Note != "" {
		fmt.Printf("Note:      %s\n", ledger.Note)
	}
	if len(ledger.Labels) > 0 {
		fmt.Printf("Labels:    ")
		for i, l := range ledger.Labels {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(l.Name)
		}
		fmt.Println()
	}
	if len(ledger.Files) > 0 {
		fmt.Printf("Files:     %d attached\n", len(ledger.Files))
	}
}

// runLedgerCreate creates a new ledger entry from flags.
func runLedgerCreate(cmd *cobra.Command, args []string) {
	client := newClient()

	amount, _ := cmd.Flags().GetFloat64("amount")
	date, _ := cmd.Flags().GetString("date")
	contactID, _ := cmd.Flags().GetUint("contact-id")
	categoryID, _ := cmd.Flags().GetUint("category-id")
	note, _ := cmd.Flags().GetString("note")

	req := &api.LedgerCreateRequest{
		Amount:     amount,
		Date:       date,
		ContactID:  contactID,
		CategoryID: categoryID,
		Note:       note,
	}

	ledger, err := client.CreateLedger(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(ledger)
		return
	}

	fmt.Printf("Created ledger entry %d (%.2f on %s)\n", ledger.ID, ledger.Amount, ledger.Date)
}

// runLedgerUpdate updates an existing ledger entry from flags.
func runLedgerUpdate(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid ledger ID")
		os.Exit(1)
	}

	req := &api.LedgerUpdateRequest{}

	if cmd.Flags().Changed("amount") {
		amount, _ := cmd.Flags().GetFloat64("amount")
		req.Amount = amount
	}
	if cmd.Flags().Changed("date") {
		date, _ := cmd.Flags().GetString("date")
		req.Date = date
	}
	if cmd.Flags().Changed("contact-id") {
		contactID, _ := cmd.Flags().GetUint("contact-id")
		req.ContactID = contactID
	}
	if cmd.Flags().Changed("category-id") {
		categoryID, _ := cmd.Flags().GetUint("category-id")
		req.CategoryID = categoryID
	}
	if cmd.Flags().Changed("note") {
		note, _ := cmd.Flags().GetString("note")
		req.Note = note
	}

	ledger, err := client.UpdateLedger(uint(id), req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(ledger)
		return
	}

	fmt.Printf("Updated ledger entry %d\n", ledger.ID)
}

// runLedgerDelete deletes a ledger entry by ID.
func runLedgerDelete(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid ledger ID")
		os.Exit(1)
	}

	if err := client.DeleteLedger(uint(id)); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted ledger entry %d\n", id)
}

// runLedgerSummary displays the ledger summary.
func runLedgerSummary(cmd *cobra.Command, args []string) {
	client := newClient()

	summary, err := client.GetLedgerSummary(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(summary)
		return
	}

	// Display years.
	if len(summary.Years) > 0 {
		fmt.Println("Years:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  YEAR\tCOUNT")
		for _, y := range summary.Years {
			fmt.Fprintf(w, "  %d\t%d\n", y.Year, y.Count)
		}
		w.Flush()
	}

	// Display categories.
	if len(summary.Categories) > 0 {
		fmt.Println("\nCategories:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  NAME\tCOUNT")
		for _, c := range summary.Categories {
			fmt.Fprintf(w, "  %s\t%d\n", c.Name, c.Count)
		}
		w.Flush()
	}

	// Display labels.
	if len(summary.Labels) > 0 {
		fmt.Println("\nLabels:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  NAME\tCOUNT")
		for _, l := range summary.Labels {
			fmt.Fprintf(w, "  %s\t%d\n", l.Name, l.Count)
		}
		w.Flush()
	}
}
