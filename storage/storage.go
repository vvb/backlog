package storage

import (
	"backlog/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	backlogDir      = "backlog"
	backlogFile     = "items.json"
	archiveFile     = "archive.json"
)

// Storage handles reading and writing backlog data
type Storage struct {
	dataDir string
}

// New creates a new Storage instance
func New() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dataDir := filepath.Join(homeDir, backlogDir)
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backlog directory: %w", err)
	}

	return &Storage{dataDir: dataDir}, nil
}

// Load reads the backlog from the JSON file
func (s *Storage) Load() (*models.Backlog, error) {
	filePath := filepath.Join(s.dataDir, backlogFile)
	
	// If file doesn't exist, return empty backlog
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &models.Backlog{Items: []models.BacklogItem{}}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read backlog file: %w", err)
	}

	var backlog models.Backlog
	if err := json.Unmarshal(data, &backlog); err != nil {
		return nil, fmt.Errorf("failed to parse backlog file: %w", err)
	}

	return &backlog, nil
}

// Save writes the backlog to the JSON file
func (s *Storage) Save(backlog *models.Backlog) error {
	filePath := filepath.Join(s.dataDir, backlogFile)
	
	data, err := json.MarshalIndent(backlog, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backlog: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backlog file: %w", err)
	}

	return nil
}

// LoadArchive reads the archived items from the JSON file
func (s *Storage) LoadArchive() (*models.Backlog, error) {
	filePath := filepath.Join(s.dataDir, archiveFile)
	
	// If file doesn't exist, return empty backlog
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &models.Backlog{Items: []models.BacklogItem{}}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read archive file: %w", err)
	}

	var backlog models.Backlog
	if err := json.Unmarshal(data, &backlog); err != nil {
		return nil, fmt.Errorf("failed to parse archive file: %w", err)
	}

	return &backlog, nil
}

// SaveArchive writes the archived items to the JSON file
func (s *Storage) SaveArchive(backlog *models.Backlog) error {
	filePath := filepath.Join(s.dataDir, archiveFile)
	
	data, err := json.MarshalIndent(backlog, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal archive: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write archive file: %w", err)
	}

	return nil
}

