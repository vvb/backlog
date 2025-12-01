package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vvb/backlog/models"
	"github.com/vvb/backlog/storage"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive completed items",
	Long:  `Move all items with 'done' status to the archive file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		// Load archive
		archive, err := store.LoadArchive()
		if err != nil {
			return err
		}

		// Separate done items from active items
		activeItems := []models.BacklogItem{}
		archivedCount := 0

		for _, item := range backlog.Items {
			if item.Status == models.StatusDone {
				archive.Items = append(archive.Items, item)
				archivedCount++
			} else {
				activeItems = append(activeItems, item)
			}
		}

		if archivedCount == 0 {
			fmt.Println("No completed items to archive")
			return nil
		}

		// Update backlog with only active items
		backlog.Items = activeItems

		// Save both files
		if err := store.Save(backlog); err != nil {
			return err
		}

		if err := store.SaveArchive(archive); err != nil {
			return err
		}

		fmt.Printf("✓ Archived %d completed item(s)\n", archivedCount)
		return nil
	},
}
