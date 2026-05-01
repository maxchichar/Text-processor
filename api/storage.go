package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	uploadsDir    = "uploads"
	processedDir  = "processed"
	historyFile   = "history.json"
	maxFileSize   = 10 * 1024 * 1024 // 10MB
)

// Storage manages file storage and history
type Storage struct {
	mu      sync.RWMutex
	history []HistoryEntry
}

// HistoryEntry represents a single processing history entry
type HistoryEntry struct {
	FileID         string            `json:"fileId"`
	Filename       string            `json:"filename"`
	OriginalSize   int64             `json:"originalSize"`
	ProcessedSize  int64             `json:"processedSize"`
	Timestamp      time.Time         `json:"timestamp"`
	Transformations ProcessOptions    `json:"transformations"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// NewStorage creates a new Storage instance
func NewStorage() (*Storage, error) {
	s := &Storage{}

	// Create directories if they don't exist
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create uploads directory: %w", err)
	}
	if err := os.MkdirAll(processedDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create processed directory: %w", err)
	}

	// Load existing history
	if err := s.loadHistory(); err != nil {
		// If history file doesn't exist, that's okay
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load history: %w", err)
		}
		s.history = []HistoryEntry{}
	}

	return s, nil
}

// generateFileID generates a unique file ID
func generateFileID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// SaveUpload saves an uploaded file and returns the file ID
func (s *Storage) SaveUpload(filename string, content io.Reader) (string, int64, error) {
	fileID := generateFileID()
	filePath := filepath.Join(uploadsDir, fileID+".txt")

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy content with size limit
	size, err := io.Copy(file, io.LimitReader(content, maxFileSize))
	if err != nil {
		os.Remove(filePath)
		return "", 0, fmt.Errorf("failed to write file: %w", err)
	}

	if size >= maxFileSize {
		os.Remove(filePath)
		return "", 0, fmt.Errorf("file size exceeds maximum limit of %d bytes", maxFileSize)
	}

	return fileID, size, nil
}

// GetUploadPath returns the path to an uploaded file
func (s *Storage) GetUploadPath(fileID string) string {
	return filepath.Join(uploadsDir, fileID+".txt")
}

// GetProcessedPath returns the path to a processed file
func (s *Storage) GetProcessedPath(fileID string) string {
	return filepath.Join(processedDir, fileID+".txt")
}

// SaveProcessed saves processed content
func (s *Storage) SaveProcessed(fileID string, content io.Reader) (int64, error) {
	filePath := s.GetProcessedPath(fileID)

	file, err := os.Create(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to create processed file: %w", err)
	}
	defer file.Close()

	size, err := io.Copy(file, content)
	if err != nil {
		return 0, fmt.Errorf("failed to write processed file: %w", err)
	}

	return size, nil
}

// DeleteFiles removes both uploaded and processed files
func (s *Storage) DeleteFiles(fileID string) error {
	uploadPath := s.GetUploadPath(fileID)
	processedPath := s.GetProcessedPath(fileID)

	var errs []error

	if err := os.Remove(uploadPath); err != nil && !os.IsNotExist(err) {
		errs = append(errs, fmt.Errorf("failed to delete upload: %w", err))
	}

	if err := os.Remove(processedPath); err != nil && !os.IsNotExist(err) {
		errs = append(errs, fmt.Errorf("failed to delete processed: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors deleting files: %v", errs)
	}

	return nil
}

// AddToHistory adds a new entry to the processing history
func (s *Storage) AddToHistory(entry HistoryEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add to beginning of history
	s.history = append([]HistoryEntry{entry}, s.history...)

	// Limit history to last 100 entries
	if len(s.history) > 100 {
		s.history = s.history[:100]
	}

	return s.saveHistory()
}

// GetHistory returns the processing history
func (s *Storage) GetHistory() []HistoryEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external modification
	history := make([]HistoryEntry, len(s.history))
	copy(history, s.history)
	return history
}

// RemoveFromHistory removes an entry from history
func (s *Storage) RemoveFromHistory(fileID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, entry := range s.history {
		if entry.FileID == fileID {
			s.history = append(s.history[:i], s.history[i+1:]...)
			return s.saveHistory()
		}
	}

	return nil
}

// loadHistory loads history from disk
func (s *Storage) loadHistory() error {
	data, err := os.ReadFile(historyFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.history)
}

// saveHistory saves history to disk
func (s *Storage) saveHistory() error {
	data, err := json.MarshalIndent(s.history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyFile, data, 0644)
}

// CleanupOldFiles removes files older than 24 hours
func (s *Storage) CleanupOldFiles() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-24 * time.Hour)
	var toDelete []string

	// Check uploads
	uploads, err := os.ReadDir(uploadsDir)
	if err == nil {
		for _, file := range uploads {
			info, err := file.Info()
			if err != nil {
				continue
			}
			if info.ModTime().Before(cutoff) {
				fileID := file.Name()[:len(file.Name())-4] // Remove .txt
				toDelete = append(toDelete, fileID)
			}
		}
	}

	// Delete old files
	for _, fileID := range toDelete {
		s.DeleteFiles(fileID)
		s.removeFromHistory(fileID)
	}

	return s.saveHistory()
}

// removeFromHistory removes an entry from history without locking
func (s *Storage) removeFromHistory(fileID string) {
	for i, entry := range s.history {
		if entry.FileID == fileID {
			s.history = append(s.history[:i], s.history[i+1:]...)
			break
		}
	}
}