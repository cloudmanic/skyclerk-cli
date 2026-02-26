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

// categoriesCmd is the parent command for category management.
var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "Manage categories",
}

// categoriesListCmd lists all categories.
var categoriesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all categories",
	Run:   runCategoriesList,
}

// categoriesGetCmd retrieves a single category by ID.
var categoriesGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a single category",
	Args:  cobra.ExactArgs(1),
	Run:   runCategoriesGet,
}

// categoriesCreateCmd creates a new category.
var categoriesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new category",
	Run:   runCategoriesCreate,
}

// categoriesUpdateCmd updates an existing category.
var categoriesUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a category",
	Args:  cobra.ExactArgs(1),
	Run:   runCategoriesUpdate,
}

// categoriesDeleteCmd deletes a category.
var categoriesDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a category",
	Args:  cobra.ExactArgs(1),
	Run:   runCategoriesDelete,
}

// init registers the categories commands and their flags.
func init() {
	// Create flags.
	categoriesCreateCmd.Flags().String("name", "", "Category name")
	categoriesCreateCmd.Flags().String("type", "", "Category type: 1=expense, 2=income")
	categoriesCreateCmd.MarkFlagRequired("name")
	categoriesCreateCmd.MarkFlagRequired("type")

	// Update flags.
	categoriesUpdateCmd.Flags().String("name", "", "Category name")
	categoriesUpdateCmd.Flags().String("type", "", "Category type: 1=expense, 2=income")

	categoriesCmd.AddCommand(categoriesListCmd)
	categoriesCmd.AddCommand(categoriesGetCmd)
	categoriesCmd.AddCommand(categoriesCreateCmd)
	categoriesCmd.AddCommand(categoriesUpdateCmd)
	categoriesCmd.AddCommand(categoriesDeleteCmd)
	rootCmd.AddCommand(categoriesCmd)
}

// runCategoriesList fetches and displays all categories.
func runCategoriesList(cmd *cobra.Command, args []string) {
	client := newClient()

	categories, err := client.GetCategories(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(categories)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tTYPE\tCOUNT")
	for _, c := range categories {
		typeLabel := capitalizeType(c.Type)
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\n", c.ID, c.Name, typeLabel, c.Count)
	}
	w.Flush()
}

// runCategoriesGet fetches and displays a single category.
func runCategoriesGet(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid category ID")
		os.Exit(1)
	}

	category, err := client.GetCategory(uint(id))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(category)
		return
	}

	fmt.Printf("ID:    %d\n", category.ID)
	fmt.Printf("Name:  %s\n", category.Name)
	fmt.Printf("Type:  %s\n", capitalizeType(category.Type))
	fmt.Printf("Count: %d\n", category.Count)
}

// runCategoriesCreate creates a new category from flags.
func runCategoriesCreate(cmd *cobra.Command, args []string) {
	client := newClient()

	name, _ := cmd.Flags().GetString("name")
	catType, _ := cmd.Flags().GetString("type")

	category, err := client.CreateCategory(&api.CategoryCreateRequest{
		Name: name,
		Type: catType,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(category)
		return
	}

	fmt.Printf("Created category %d: %s\n", category.ID, category.Name)
}

// runCategoriesUpdate updates an existing category from flags.
func runCategoriesUpdate(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid category ID")
		os.Exit(1)
	}

	req := &api.CategoryUpdateRequest{}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		req.Name = name
	}
	if cmd.Flags().Changed("type") {
		catType, _ := cmd.Flags().GetString("type")
		req.Type = catType
	}

	category, err := client.UpdateCategory(uint(id), req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(category)
		return
	}

	fmt.Printf("Updated category %d: %s\n", category.ID, category.Name)
}

// runCategoriesDelete deletes a category by ID.
func runCategoriesDelete(cmd *cobra.Command, args []string) {
	client := newClient()

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid category ID")
		os.Exit(1)
	}

	if err := client.DeleteCategory(uint(id)); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted category %d\n", id)
}

// capitalizeType returns the category type string with the first letter capitalized.
func capitalizeType(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
