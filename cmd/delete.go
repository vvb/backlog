package cmd

import (
	"backlog/models"
	"backlog/storage"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a backlog item",
	Long:  `Delete a backlog item by its ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

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

		// Find and remove item
		found := false
		newItems := []models.BacklogItem{}
		var deletedTitle string

		for _, item := range backlog.Items {
			if strings.HasPrefix(item.ID, id) {
				found = true
				deletedTitle = item.Title
			} else {
				newItems = append(newItems, item)
			}
		}

		if !found {
			return fmt.Errorf("item with ID %s not found", id)
		}

		backlog.Items = newItems

		// Save
		if err := store.Save(backlog); err != nil {
			return err
		}

		fmt.Printf("✓ Deleted backlog item: %s\n", deletedTitle)
		return nil
	},
}
