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

// labelsCmd is the parent command for label management.
var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage labels",
}

// labelsListCmd lists all labels.
var labelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all labels",
	Run:   runLabelsList,
}

// labelsGetCmd retrieves a single label by ID.
var labelsGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a single label",
	Args:  cobra.ExactArgs(1),
	Run:   runLabelsGet,
}

// labelsCreateCmd creates a new label.
var labelsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new label",
	Run:   runLabelsCreate,
}

// labelsUpdateCmd updates an existing label.
var labelsUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a label",
	Args:  cobra.ExactArgs(1),
	Run:   runLabelsUpdate,
}

// labelsDeleteCmd deletes a label.
var labelsDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a label",
	Args:  cobra.ExactArgs(1),
	Run:   runLabelsDelete,
}

// init registers the labels commands and their flags.
func init() {
	// Create flags.
	labelsCreateCmd.Flags().String("name", "", "Label name")
	labelsCreateCmd.MarkFlagRequired("name")

	// Update flags.
	labelsUpdateCmd.Flags().String("name", "", "Label name")

	labelsCmd.AddCommand(labelsListCmd)
	labelsCmd.AddCommand(labelsGetCmd)
	labelsCmd.AddCommand(labelsCreateCmd)
	labelsCmd.AddCommand(labelsUpdateCmd)
	labelsCmd.AddCommand(labelsDeleteCmd)
	rootCmd.AddCommand(labelsCmd)
}

// runLabelsList fetches and displays all labels.
func runLabelsList(cmd *cobra.Command, args []string) {
	client := newClient()

	labels, err := client.GetLabels(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(labels)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCOUNT")
	for _, l := range labels {
		fmt.Fprintf(w, "%d\t%s\t%d\n", l.ID, l.Name, l.Count)
	}
	w.Flush()
}

// runLabelsGet fetches and displays a single label.
func runLabelsGet(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid label ID")
		os.Exit(1)
	}

	label, err := client.GetLabel(uint(id))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(label)
		return
	}

	fmt.Printf("ID:    %d\n", label.ID)
	fmt.Printf("Name:  %s\n", label.Name)
	fmt.Printf("Count: %d\n", label.Count)
}

// runLabelsCreate creates a new label from flags.
func runLabelsCreate(cmd *cobra.Command, args []string) {
	client := newClient()

	name, _ := cmd.Flags().GetString("name")

	label, err := client.CreateLabel(&api.LabelCreateRequest{Name: name})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(label)
		return
	}

	fmt.Printf("Created label %d: %s\n", label.ID, label.Name)
}

// runLabelsUpdate updates an existing label from flags.
func runLabelsUpdate(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid label ID")
		os.Exit(1)
	}

	name, _ := cmd.Flags().GetString("name")

	label, err := client.UpdateLabel(uint(id), &api.LabelUpdateRequest{Name: name})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(label)
		return
	}

	fmt.Printf("Updated label %d: %s\n", label.ID, label.Name)
}

// runLabelsDelete deletes a label by ID.
func runLabelsDelete(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid label ID")
		os.Exit(1)
	}

	if err := client.DeleteLabel(uint(id)); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted label %d\n", id)
}
