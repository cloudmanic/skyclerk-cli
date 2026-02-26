// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

// outputFormat controls whether output is displayed as a table or JSON.
var outputFormat string

// accountOverride allows overriding the default account ID for a single command.
var accountOverride uint

// rootCmd is the base command for the Skyclerk CLI.
var rootCmd = &cobra.Command{
	Use:   "skyclerk",
	Short: "Skyclerk CLI - manage your bookkeeping from the terminal",
	Long:  "A command-line interface for the Skyclerk bookkeeping API.",
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// init registers global persistent flags available to all commands.
func init() {
	rootCmd.PersistentFlags().StringVar(&outputFormat, "output", "table", "Output format: table or json")
	rootCmd.PersistentFlags().UintVar(&accountOverride, "account", 0, "Override the default account ID")
}
