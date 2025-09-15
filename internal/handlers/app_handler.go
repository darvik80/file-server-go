package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"file-server-go/internal/database"
	"file-server-go/internal/middleware"
	"file-server-go/internal/services"
)

// AppHandler handles application-related requests
type AppHandler struct {
	db          *database.DB
	appService  *services.ApplicationService
	fileService *services.FileService
	uploadDir   string
}

// NewAppHandler creates new application handler
func NewAppHandler(db *database.DB, uploadDir string) *AppHandler {
	return &AppHandler{
		db:          db,
		appService:  services.NewApplicationService(db),
		fileService: services.NewFileService(db),
		uploadDir:   uploadDir,
	}
}

// RegisterAppRequest represents application registration request
type RegisterAppRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RegisterAppResponse represents successful application registration response
type RegisterAppResponse struct {
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Name            string `json:"name"`
	Description     string `json:"description"`
}

// RegisterApp handles application registration
func (h *AppHandler) RegisterApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		sendError(w, "Application name is required", http.StatusBadRequest)
		return
	}

	// Register application
	app, err := h.appService.RegisterApplication(req.Name, req.Description)
	if err != nil {
		sendError(w, "Failed to register application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	response := RegisterAppResponse{
		AccessKeyID:     app.AccessKeyID,
		AccessKeySecret: app.AccessKeySecret,
		Name:            app.Name,
		Description:     app.Description,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		sendError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UploadRawFileResponse represents successful raw file upload response
type UploadRawFileResponse struct {
	URL       string `json:"url"`
	Filename  string `json:"filename"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

// UploadRawFile handles raw file upload for applications
func (h *AppHandler) UploadRawFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Check if application is authenticated
	app := middleware.GetAppFromContext(r.Context())
	if app == nil {
		sendError(w, "Application authentication required", http.StatusUnauthorized)
		return
	}

	// Get filename from URL path or query parameter
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		// Try to get filename from URL path
		filename = r.URL.Path[len("/upload/raw/"):]
	}

	if filename == "" {
		sendError(w, "Filename not specified", http.StatusBadRequest)
		return
	}

	// Get TTL from query parameter (in seconds)
	var ttl time.Duration
	if ttlStr := r.URL.Query().Get("ttl"); ttlStr != "" {
		ttlSeconds, err := strconv.ParseInt(ttlStr, 10, 64)
		if err != nil {
			sendError(w, "Invalid TTL value", http.StatusBadRequest)
			return
		}
		ttl = time.Duration(ttlSeconds) * time.Second
	}

	// Create file path
	filePath := filepath.Join(h.uploadDir, filename)

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		sendError(w, "Failed to create directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create file
	dst, err := os.Create(filePath)
	if err != nil {
		sendError(w, "Error creating file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			// Log error but don't fail the request
		}
	}()

	// Copy request body to file
	_, err = io.Copy(dst, r.Body)
	if err != nil {
		sendError(w, "Error saving file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate URL for downloading the file
	url := "/download/" + filename

	// Create file record in database
	fileRecord, err := h.fileService.CreateFileRecord(app.ID, filename, filePath, url, ttl)
	if err != nil {
		sendError(w, "Failed to create file record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := UploadRawFileResponse{
		URL:      url,
		Filename: filename,
	}

	// Add expiration time to response if TTL was specified
	if ttl > 0 {
		response.ExpiresAt = fileRecord.ExpirationTime.Format(time.RFC3339)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		sendError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
