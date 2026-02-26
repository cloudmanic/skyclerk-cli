// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// reportsCmd is the parent command for reports.
var reportsCmd = &cobra.Command{
	Use:   "reports",
	Short: "Generate financial reports",
}

// reportsPnlCmd generates a P&L report for a date range.
var reportsPnlCmd = &cobra.Command{
	Use:   "pnl",
	Short: "Profit and loss report",
	Run:   runReportsPnl,
}

// reportsPnlCurrentCmd generates the current year P&L report.
var reportsPnlCurrentCmd = &cobra.Command{
	Use:   "pnl-current",
	Short: "Current year profit and loss report",
	Run:   runReportsPnlCurrent,
}

// reportsPnlByCategoryCmd generates a P&L report by category.
var reportsPnlByCategoryCmd = &cobra.Command{
	Use:   "pnl-by-category",
	Short: "Profit and loss by category",
	Run:   runReportsPnlByCategory,
}

// reportsPnlByLabelCmd generates a P&L report by label.
var reportsPnlByLabelCmd = &cobra.Command{
	Use:   "pnl-by-label",
	Short: "Profit and loss by label",
	Run:   runReportsPnlByLabel,
}

// reportsIncomeByContactCmd generates an income by contact report.
var reportsIncomeByContactCmd = &cobra.Command{
	Use:   "income-by-contact",
	Short: "Income breakdown by contact",
	Run:   runReportsIncomeByContact,
}

// reportsExpensesByContactCmd generates an expenses by contact report.
var reportsExpensesByContactCmd = &cobra.Command{
	Use:   "expenses-by-contact",
	Short: "Expenses breakdown by contact",
	Run:   runReportsExpensesByContact,
}

// init registers the reports commands and their flags.
func init() {
	// P&L date range flags.
	reportsPnlCmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
	reportsPnlCmd.Flags().String("end", "", "End date (YYYY-MM-DD)")

	reportsPnlByCategoryCmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
	reportsPnlByCategoryCmd.Flags().String("end", "", "End date (YYYY-MM-DD)")

	reportsPnlByLabelCmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
	reportsPnlByLabelCmd.Flags().String("end", "", "End date (YYYY-MM-DD)")

	reportsIncomeByContactCmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
	reportsIncomeByContactCmd.Flags().String("end", "", "End date (YYYY-MM-DD)")

	reportsExpensesByContactCmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
	reportsExpensesByContactCmd.Flags().String("end", "", "End date (YYYY-MM-DD)")

	reportsCmd.AddCommand(reportsPnlCmd)
	reportsCmd.AddCommand(reportsPnlCurrentCmd)
	reportsCmd.AddCommand(reportsPnlByCategoryCmd)
	reportsCmd.AddCommand(reportsPnlByLabelCmd)
	reportsCmd.AddCommand(reportsIncomeByContactCmd)
	reportsCmd.AddCommand(reportsExpensesByContactCmd)
	rootCmd.AddCommand(reportsCmd)
}

// getDateParams extracts start and end date params from flags.
func getDateParams(cmd *cobra.Command) map[string]string {
	params := map[string]string{}

	start, _ := cmd.Flags().GetString("start")
	end, _ := cmd.Flags().GetString("end")

	if start != "" {
		params["start"] = start
	}
	if end != "" {
		params["end"] = end
	}

	return params
}

// runReportsPnl generates and displays a P&L report.
func runReportsPnl(cmd *cobra.Command, args []string) {
	client := newClient()

	report, err := client.GetPnlReport(getDateParams(cmd))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(report)
		return
	}

	fmt.Printf("Income:  %.2f\n", report.Income)
	fmt.Printf("Expense: %.2f\n", report.Expense)
	fmt.Printf("Profit:  %.2f\n", report.Profit)

	if len(report.Breakdown) > 0 {
		fmt.Println("\nBreakdown:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  NAME\tAMOUNT")
		for _, b := range report.Breakdown {
			fmt.Fprintf(w, "  %s\t%.2f\n", b.Name, b.Amount)
		}
		w.Flush()
	}
}

// runReportsPnlCurrent generates and displays the current year P&L.
func runReportsPnlCurrent(cmd *cobra.Command, args []string) {
	client := newClient()

	report, err := client.GetPnlCurrent()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(report)
		return
	}

	fmt.Printf("Income:  %.2f\n", report.Income)
	fmt.Printf("Expense: %.2f\n", report.Expense)
	fmt.Printf("Profit:  %.2f\n", report.Profit)
}

// runReportsPnlByCategory generates and displays a P&L by category.
func runReportsPnlByCategory(cmd *cobra.Command, args []string) {
	client := newClient()

	report, err := client.GetPnlByCategory(getDateParams(cmd))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(report)
		return
	}

	fmt.Printf("Income:  %.2f\n", report.Income)
	fmt.Printf("Expense: %.2f\n", report.Expense)
	fmt.Printf("Profit:  %.2f\n", report.Profit)

	if len(report.Breakdown) > 0 {
		fmt.Println("\nBy Category:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  CATEGORY\tAMOUNT")
		for _, b := range report.Breakdown {
			fmt.Fprintf(w, "  %s\t%.2f\n", b.Name, b.Amount)
		}
		w.Flush()
	}
}

// runReportsPnlByLabel generates and displays a P&L by label.
func runReportsPnlByLabel(cmd *cobra.Command, args []string) {
	client := newClient()

	report, err := client.GetPnlByLabel(getDateParams(cmd))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(report)
		return
	}

	fmt.Printf("Income:  %.2f\n", report.Income)
	fmt.Printf("Expense: %.2f\n", report.Expense)
	fmt.Printf("Profit:  %.2f\n", report.Profit)

	if len(report.Breakdown) > 0 {
		fmt.Println("\nBy Label:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  LABEL\tAMOUNT")
		for _, b := range report.Breakdown {
			fmt.Fprintf(w, "  %s\t%.2f\n", b.Name, b.Amount)
		}
		w.Flush()
	}
}

// runReportsIncomeByContact generates and displays income by contact.
func runReportsIncomeByContact(cmd *cobra.Command, args []string) {
	client := newClient()

	report, err := client.GetIncomeByContact(getDateParams(cmd))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(report)
		return
	}

	fmt.Printf("Total Income: %.2f\n", report.Income)

	if len(report.Breakdown) > 0 {
		fmt.Println("\nBy Contact:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  CONTACT\tAMOUNT")
		for _, b := range report.Breakdown {
			fmt.Fprintf(w, "  %s\t%.2f\n", b.Name, b.Amount)
		}
		w.Flush()
	}
}

// runReportsExpensesByContact generates and displays expenses by contact.
func runReportsExpensesByContact(cmd *cobra.Command, args []string) {
	client := newClient()

	report, err := client.GetExpensesByContact(getDateParams(cmd))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(report)
		return
	}

	fmt.Printf("Total Expenses: %.2f\n", report.Expense)

	if len(report.Breakdown) > 0 {
		fmt.Println("\nBy Contact:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  CONTACT\tAMOUNT")
		for _, b := range report.Breakdown {
			fmt.Fprintf(w, "  %s\t%.2f\n", b.Name, b.Amount)
		}
		w.Flush()
	}
}
