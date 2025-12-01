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
	addDesc    string
	addDueDate string
	addTags    string
)

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new backlog item",
	Long:  `Add a new backlog item with title, description, due date, and tags.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]

		// Validate due date format if provided
		if addDueDate != "" {
			if !isValidDateFormat(addDueDate) {
				return fmt.Errorf("invalid date format. Use DD-MM-YYYY")
			}
		}

		// Parse tags
		var tags []string
		if addTags != "" {
			tags = strings.Split(addTags, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		// Create storage
		store, err := storage.New()
		if err != nil {
			return err
		}

		// Load existing backlog
		backlog, err := store.Load()
		if err != nil {
			return err
		}

		// Generate ID
		id := generateID()

		// Create new item
		item := models.BacklogItem{
			ID:          id,
			Title:       title,
			Description: addDesc,
			DueDate:     addDueDate,
			Tags:        tags,
			Status:      models.StatusTodo,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Add to backlog
		backlog.Items = append(backlog.Items, item)

		// Save
		if err := store.Save(backlog); err != nil {
			return err
		}

		fmt.Printf("✓ Added backlog item: %s (ID: %s)\n", title, id)
		return nil
	},
}

func init() {
	addCmd.Flags().StringVar(&addDesc, "desc", "", "Description of the backlog item")
	addCmd.Flags().StringVar(&addDueDate, "due", "", "Due date in DD-MM-YYYY format")
	addCmd.Flags().StringVar(&addTags, "tags", "", "Comma-separated tags")
}

// isValidDateFormat checks if the date is in DD-MM-YYYY format
func isValidDateFormat(date string) bool {
	_, err := time.Parse("02-01-2006", date)
	return err == nil
}

// generateID generates a simple unique ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
