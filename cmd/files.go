// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// filesCmd is the parent command for file management.
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Manage files and receipts",
}

// filesUploadCmd uploads a file to the current account.
var filesUploadCmd = &cobra.Command{
	Use:   "upload [file-path]",
	Short: "Upload a file or receipt",
	Args:  cobra.ExactArgs(1),
	Run:   runFilesUpload,
}

// init registers the files commands and their flags.
func init() {
	filesUploadCmd.Flags().String("ledger-id", "", "Associate file with a ledger entry")

	filesCmd.AddCommand(filesUploadCmd)
	rootCmd.AddCommand(filesCmd)
}

// runFilesUpload uploads a file to the Skyclerk API.
func runFilesUpload(cmd *cobra.Command, args []string) {
	client := newClient()

	filePath := args[0]
	ledgerID, _ := cmd.Flags().GetString("ledger-id")

	// Verify the file exists before uploading.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: file not found: %s\n", filePath)
		os.Exit(1)
	}

	file, err := client.UploadFile(filePath, ledgerID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		printJSON(file)
		return
	}

	fmt.Printf("Uploaded file %d: %s (%s, %d bytes)\n", file.ID, file.Name, file.Type, file.Size)
}
