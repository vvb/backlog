package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "backlog",
	Short: "Backlog - A terminal application for managing backlog items",
	Long:  `Backlog is a CLI tool for creating and managing backlog items with a Kanban-style board view.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(archiveCmd)
}
