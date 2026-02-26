// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// activitiesCmd lists recent account activities.
var activitiesCmd = &cobra.Command{
	Use:   "activities",
	Short: "List recent account activities",
	Run:   runActivities,
}

// init registers the activities command and its flags.
func init() {
	activitiesCmd.Flags().String("limit", "100", "Number of activities to return")
	activitiesCmd.Flags().String("order", "id", "Sort field")
	activitiesCmd.Flags().String("sort", "DESC", "Sort direction (ASC or DESC)")

	rootCmd.AddCommand(activitiesCmd)
}

// runActivities fetches and displays recent account activities.
func runActivities(cmd *cobra.Command, args []string) {
	client := newClient()

	limit, _ := cmd.Flags().GetString("limit")
	order, _ := cmd.Flags().GetString("order")
	sort, _ := cmd.Flags().GetString("sort")

	params := map[string]string{
		"limit": limit,
		"order": order,
		"sort":  sort,
	}

	activities, err := client.GetActivities(params)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(activities)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tACTION\tMESSAGE\tDATE")
	for _, a := range activities {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", a.ID, a.Action, a.Message, a.CreatedAt)
	}
	w.Flush()
}
