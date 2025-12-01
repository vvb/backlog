package cmd

import (
	"backlog/models"
	"backlog/storage"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search for backlog items",
	Long:  `Search for backlog items by keyword in title, description, or tags.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyword := strings.ToLower(args[0])

		// Create storage
		store, err := storage.New()
		if err != nil {
			return err
		}

		// Load backlog
		backlog, err := store.Load()
		if err != nil {
			return err
		}

		// Search items
		matches := []models.BacklogItem{}
		for _, item := range backlog.Items {
			if matchesKeyword(item, keyword) {
				matches = append(matches, item)
			}
		}

		// Display results
		if len(matches) == 0 {
			fmt.Printf("No items found matching '%s'\n", args[0])
			return nil
		}

		fmt.Printf("\nFound %d item(s) matching '%s':\n\n", len(matches), args[0])
		for _, item := range matches {
			displayItem(item)
			fmt.Println()
		}

		return nil
	},
}

func matchesKeyword(item models.BacklogItem, keyword string) bool {
	// Check title
	if strings.Contains(strings.ToLower(item.Title), keyword) {
		return true
	}

	// Check description
	if strings.Contains(strings.ToLower(item.Description), keyword) {
		return true
	}

	// Check tags
	for _, tag := range item.Tags {
		if strings.Contains(strings.ToLower(tag), keyword) {
			return true
		}
	}

	return false
}

func displayItem(item models.BacklogItem) {
	fmt.Printf("ID: %s\n", truncateID(item.ID))
	fmt.Printf("Title: %s\n", item.Title)
	if item.Description != "" {
		fmt.Printf("Description: %s\n", item.Description)
	}
	fmt.Printf("Status: %s\n", item.Status)
	if item.DueDate != "" {
		fmt.Printf("Due Date: %s\n", item.DueDate)
	}
	if len(item.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(item.Tags, ", "))
	}
	fmt.Printf("Created: %s\n", item.CreatedAt.Format("02-01-2006 15:04"))
	fmt.Printf("Updated: %s\n", item.UpdatedAt.Format("02-01-2006 15:04"))
}

