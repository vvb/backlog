package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/vvb/backlog/models"
	"github.com/vvb/backlog/storage"
)

var interactive bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all backlog items in Kanban board view",
	Long:  `Display all backlog items organized in a Kanban-style board with columns for todo, in-progress, and done.`,
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

		// Interactive mode
		if interactive {
			p := tea.NewProgram(initialModel(backlog, store))
			if _, err := p.Run(); err != nil {
				return err
			}
			return nil
		}

		// Display Kanban board
		displayKanbanBoard(backlog)
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode with keyboard navigation")
}

func displayKanbanBoard(backlog *models.Backlog) {
	// Organize items by status
	todoItems := []models.BacklogItem{}
	inProgressItems := []models.BacklogItem{}
	doneItems := []models.BacklogItem{}

	for _, item := range backlog.Items {
		switch item.Status {
		case models.StatusTodo:
			todoItems = append(todoItems, item)
		case models.StatusInProgress:
			inProgressItems = append(inProgressItems, item)
		case models.StatusDone:
			doneItems = append(doneItems, item)
		}
	}

	// Calculate column width
	const colWidth = 35

	// Print header
	fmt.Println()
	fmt.Println(strings.Repeat("=", colWidth*3+4))
	fmt.Printf("%-*s | %-*s | %-*s\n", colWidth, "TODO", colWidth, "IN PROGRESS", colWidth, "DONE")
	fmt.Println(strings.Repeat("=", colWidth*3+4))

	// Find max rows needed
	maxRows := max(len(todoItems), len(inProgressItems), len(doneItems))

	// Print rows
	for i := 0; i < maxRows; i++ {
		todoCell := formatCell(todoItems, i, colWidth)
		inProgressCell := formatCell(inProgressItems, i, colWidth)
		doneCell := formatCell(doneItems, i, colWidth)

		fmt.Printf("%s | %s | %s\n", todoCell, inProgressCell, doneCell)

		// Add separator between items
		if i < maxRows-1 {
			fmt.Printf("%s | %s | %s\n",
				strings.Repeat("-", colWidth),
				strings.Repeat("-", colWidth),
				strings.Repeat("-", colWidth))
		}
	}

	fmt.Println(strings.Repeat("=", colWidth*3+4))
	fmt.Printf("\nTotal: %d items (%d todo, %d in-progress, %d done)\n\n",
		len(backlog.Items), len(todoItems), len(inProgressItems), len(doneItems))
}

func formatCell(items []models.BacklogItem, index int, width int) string {
	if index >= len(items) {
		return strings.Repeat(" ", width)
	}

	item := items[index]

	// Format: [ID] Title
	// Tags: tag1, tag2
	// Due: DD-MM-YYYY

	line1 := fmt.Sprintf("[%s] %s", truncateID(item.ID), item.Title)
	line1 = truncate(line1, width)

	lines := []string{line1}

	if len(item.Tags) > 0 {
		tagsStr := "Tags: " + strings.Join(item.Tags, ", ")
		lines = append(lines, truncate(tagsStr, width))
	}

	if item.DueDate != "" {
		dueStr := "Due: " + item.DueDate
		lines = append(lines, truncate(dueStr, width))
	}

	// Return the first line (simplified for single-line display)
	return truncate(line1, width)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s + strings.Repeat(" ", maxLen-len(s))
	}
	return s[:maxLen-3] + "..."
}

func truncateID(id string) string {
	if len(id) <= 6 {
		return id
	}
	return id[:6]
}

func max(a, b, c int) int {
	result := a
	if b > result {
		result = b
	}
	if c > result {
		result = c
	}
	return result
}
