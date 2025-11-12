/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "csvtk",
	Short: "A powerful CSV utility toolkit",
	Long: `csvtk is a command-line tool for working with CSV files.

It provides utilities for viewing, editing, validating, and transforming CSV data,
including features like counting rows/columns, moving data, filtering, and more.

When called with just a filename, it will open an interactive viewer.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {

			viewCommand.Run(cmd, args)
		} else {
			cmd.Help()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
