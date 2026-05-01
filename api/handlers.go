package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Server represents the API server
type Server struct {
	storage  *Storage
	processor *Processor
}

// NewServer creates a new API server
func NewServer() (*Server, error) {
	storage, err := NewStorage()
	if err != nil {
		return nil, err
	}

	return &Server{
		storage:  storage,
		processor: NewProcessor(),
	}, nil
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// UploadResponse represents the response from file upload
type UploadResponse struct {
	FileID   string `json:"fileId"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// ProcessResponse represents the response from file processing
type ProcessResponse struct {
	FileID string `json:"fileId"`
	Status string `json:"status"`
}

// PreviewResponse represents the response from preview request
type PreviewResponse struct {
	Original  string `json:"original"`
	Processed string `json:"processed"`
}

// HistoryResponse represents the response from history request
type HistoryResponse struct {
	History []HistoryEntry `json:"history"`
}

// DeleteResponse represents the response from delete request
type DeleteResponse struct {
	Success bool `json:"success"`
}

// sendJSON sends a JSON response
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// sendError sends an error response
func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, ErrorResponse{Error: message})
}

// setCORS sets CORS headers
func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// handleUpload handles file upload requests
func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	setCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse multipart form (max 32MB)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		sendError(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// Validate file extension
	if !strings.HasSuffix(header.Filename, ".txt") {
		sendError(w, http.StatusBadRequest, "Only .txt files are allowed")
		return
	}

	// Save the file
	fileID, size, err := s.storage.SaveUpload(header.Filename, file)
	if err != nil {
		log.Printf("Error saving upload: %v", err)
		sendError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}

	sendJSON(w, http.StatusOK, UploadResponse{
		FileID:   fileID,
		Filename: header.Filename,
		Size:     size,
	})
}

// handleProcess handles file processing requests
func (s *Server) handleProcess(w http.ResponseWriter, r *http.Request) {
	setCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract file ID from URL
	fileID := strings.TrimPrefix(r.URL.Path, "/api/process/")
	if fileID == "" {
		sendError(w, http.StatusBadRequest, "File ID required")
		return
	}

	// Parse request body
	var options ProcessOptions
	if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get paths
	inputPath := s.storage.GetUploadPath(fileID)
	outputPath := s.storage.GetProcessedPath(fileID)

	// Check if input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		sendError(w, http.StatusNotFound, "File not found")
		return
	}

	// Process the file
	if err := s.processor.ProcessFile(inputPath, outputPath, options); err != nil {
		log.Printf("Error processing file: %v", err)
		sendError(w, http.StatusInternalServerError, "Failed to process file")
		return
	}

	// Get file info for history
	inputInfo, _ := os.Stat(inputPath)
	outputInfo, _ := os.Stat(outputPath)

	// Add to history
	entry := HistoryEntry{
		FileID:         fileID,
		Filename:       filepath.Base(inputPath),
		OriginalSize:   inputInfo.Size(),
		ProcessedSize:  outputInfo.Size(),
		Timestamp:      inputInfo.ModTime(),
		Transformations: options,
	}
	s.storage.AddToHistory(entry)

	sendJSON(w, http.StatusOK, ProcessResponse{
		FileID: fileID,
		Status: "processed",
	})
}

// handlePreview handles preview requests
func (s *Server) handlePreview(w http.ResponseWriter, r *http.Request) {
	setCORS(w)

	if r.Method != "GET" {
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract file ID from URL
	fileID := strings.TrimPrefix(r.URL.Path, "/api/preview/")
	if fileID == "" {
		sendError(w, http.StatusBadRequest, "File ID required")
		return
	}

	// Get paths
	inputPath := s.storage.GetUploadPath(fileID)
	outputPath := s.storage.GetProcessedPath(fileID)

	// Read original file
	original, err := os.ReadFile(inputPath)
	if err != nil {
		sendError(w, http.StatusNotFound, "Original file not found")
		return
	}

	// Read processed file
	processed, err := os.ReadFile(outputPath)
	if err != nil {
		sendError(w, http.StatusNotFound, "Processed file not found")
		return
	}

	sendJSON(w, http.StatusOK, PreviewResponse{
		Original:  string(original),
		Processed: string(processed),
	})
}

// handleDownload handles file download requests
func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	setCORS(w)

	if r.Method != "GET" {
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract file ID from URL
	fileID := strings.TrimPrefix(r.URL.Path, "/api/download/")
	if fileID == "" {
		sendError(w, http.StatusBadRequest, "File ID required")
		return
	}

	// Get processed file path
	filePath := s.storage.GetProcessedPath(fileID)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		sendError(w, http.StatusNotFound, "Processed file not found")
		return
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to open file")
		return
	}
	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to get file info")
		return
	}

	// Set headers for download
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=processed-%s.txt", fileID))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

	// Copy file to response
	io.Copy(w, file)
}

// handleHistory handles history requests
func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	setCORS(w)

	if r.Method != "GET" {
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	history := s.storage.GetHistory()
	sendJSON(w, http.StatusOK, HistoryResponse{
		History: history,
	})
}

// handleDelete handles file deletion requests
func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	setCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "DELETE" {
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract file ID from URL
	fileID := strings.TrimPrefix(r.URL.Path, "/api/files/")
	if fileID == "" {
		sendError(w, http.StatusBadRequest, "File ID required")
		return
	}

	// Delete files
	if err := s.storage.DeleteFiles(fileID); err != nil {
		log.Printf("Error deleting files: %v", err)
		sendError(w, http.StatusInternalServerError, "Failed to delete files")
		return
	}

	// Remove from history
	s.storage.RemoveFromHistory(fileID)

	sendJSON(w, http.StatusOK, DeleteResponse{
		Success: true,
	})
}

// RegisterRoutes registers all API routes
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/upload", s.handleUpload)
	mux.HandleFunc("/api/process/", s.handleProcess)
	mux.HandleFunc("/api/preview/", s.handlePreview)
	mux.HandleFunc("/api/download/", s.handleDownload)
	mux.HandleFunc("/api/history", s.handleHistory)
	mux.HandleFunc("/api/files/", s.handleDelete)
}