package models

import (
	"time"
)

// Status represents the status of a backlog item
type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

// BacklogItem represents a single backlog item
type BacklogItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"` // Format: DD-MM-YYYY
	Tags        []string  `json:"tags"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Backlog represents the collection of all backlog items
type Backlog struct {
	Items []BacklogItem `json:"items"`
}

// ValidStatus checks if a status string is valid
func ValidStatus(s string) bool {
	return s == string(StatusTodo) || s == string(StatusInProgress) || s == string(StatusDone)
}

