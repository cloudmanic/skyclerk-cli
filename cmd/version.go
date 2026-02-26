// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd displays the CLI version.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("skyclerk version %s\n", Version)
	},
}

// init registers the version command.
func init() {
	rootCmd.AddCommand(versionCmd)
}
