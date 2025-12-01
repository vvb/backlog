package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vvb/backlog/models"
	"github.com/vvb/backlog/storage"
)

var (
	updateTitle  string
	updateDesc   string
	updateDue    string
	updateTags   string
	updateStatus string
)

var updateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a backlog item",
	Long:  `Update a backlog item's title, description, due date, tags, or status.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		// Validate status if provided
		if updateStatus != "" && !models.ValidStatus(updateStatus) {
			return fmt.Errorf("invalid status. Use: todo, in-progress, or done")
		}

		// Validate due date format if provided
		if updateDue != "" && !isValidDateFormat(updateDue) {
			return fmt.Errorf("invalid date format. Use DD-MM-YYYY")
		}

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

		// Find item
		found := false
		for i := range backlog.Items {
			if strings.HasPrefix(backlog.Items[i].ID, id) {
				found = true

				// Update fields
				if updateTitle != "" {
					backlog.Items[i].Title = updateTitle
				}
				if updateDesc != "" {
					backlog.Items[i].Description = updateDesc
				}
				if updateDue != "" {
					backlog.Items[i].DueDate = updateDue
				}
				if updateTags != "" {
					tags := strings.Split(updateTags, ",")
					for j := range tags {
						tags[j] = strings.TrimSpace(tags[j])
					}
					backlog.Items[i].Tags = tags
				}
				if updateStatus != "" {
					backlog.Items[i].Status = models.Status(updateStatus)
				}

				backlog.Items[i].UpdatedAt = time.Now()

				// Save
				if err := store.Save(backlog); err != nil {
					return err
				}

				fmt.Printf("✓ Updated backlog item: %s\n", backlog.Items[i].Title)
				break
			}
		}

		if !found {
			return fmt.Errorf("item with ID %s not found", id)
		}

		return nil
	},
}

func init() {
	updateCmd.Flags().StringVar(&updateTitle, "title", "", "New title")
	updateCmd.Flags().StringVar(&updateDesc, "desc", "", "New description")
	updateCmd.Flags().StringVar(&updateDue, "due", "", "New due date in DD-MM-YYYY format")
	updateCmd.Flags().StringVar(&updateTags, "tags", "", "New comma-separated tags")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "New status (todo, in-progress, done)")
}
