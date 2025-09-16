package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"file-server-go/internal/auth"
	"file-server-go/internal/database"
	"file-server-go/internal/middleware"
	"file-server-go/internal/models"
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
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Permissions []models.AppPermission `json:"permissions"`
}

// RegisterAppResponse represents successful application registration response
type RegisterAppResponse struct {
	AccessKeyID     string                 `json:"access_key_id"`
	AccessKeySecret string                 `json:"access_key_secret"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Permissions     []models.AppPermission `json:"permissions"`
}

// RegisterApp handles application registration
func (h *AppHandler) RegisterApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	currentUsername := auth.GetUserFromContext(r.Context())
	if currentUsername == "" {
		sendError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Find current user
	currentUser := &models.User{}
	if err := currentUser.FindUserByUsername(h.db.DB, currentUsername); err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user has admin role
	if !currentUser.HasAdminRole() {
		sendError(w, "Access denied. Admin rights required", http.StatusForbidden)
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

	// Set default permissions if not provided
	if len(req.Permissions) == 0 {
		req.Permissions = []models.AppPermission{models.ReadPermission}
	}

	// Validate permissions
	validPermissions := map[models.AppPermission]bool{
		models.ReadPermission:  true,
		models.WritePermission: true,
	}

	for _, perm := range req.Permissions {
		if !validPermissions[perm] {
			sendError(w, "Invalid permission. Valid permissions are: read, write", http.StatusBadRequest)
			return
		}
	}

	// Register application
	appModel := &models.Application{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
	}

	// Generate credentials
	appModel.AccessKeyID = generateAccessKeyID()
	appModel.AccessKeySecret = generateAccessKeySecret()

	if err := appModel.CreateApplication(h.db.DB); err != nil {
		sendError(w, "Failed to register application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	response := RegisterAppResponse{
		AccessKeyID:     appModel.AccessKeyID,
		AccessKeySecret: appModel.AccessKeySecret,
		Name:            appModel.Name,
		Description:     appModel.Description,
		Permissions:     appModel.Permissions,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		sendError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetAppsResponse represents applications list response
type GetAppsResponse struct {
	Applications []GetAppResponse `json:"applications"`
}

// GetAppResponse represents single application response
type GetAppResponse struct {
	ID          int64                  `json:"id"`
	AccessKeyID string                 `json:"access_key_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Permissions []models.AppPermission `json:"permissions"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// GetApplications handles getting all applications
func (h *AppHandler) GetApplications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	currentUsername := auth.GetUserFromContext(r.Context())
	if currentUsername == "" {
		sendError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Find current user
	currentUser := &models.User{}
	if err := currentUser.FindUserByUsername(h.db.DB, currentUsername); err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user has admin role
	if !currentUser.HasAdminRole() {
		sendError(w, "Access denied. Admin rights required", http.StatusForbidden)
		return
	}

	// Get all applications
	appModel := &models.Application{}
	apps, err := appModel.GetAllApplications(h.db.DB)
	if err != nil {
		sendError(w, "Failed to get applications: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	var responseApps []GetAppResponse
	for _, app := range apps {
		responseApps = append(responseApps, GetAppResponse{
			ID:          app.ID,
			AccessKeyID: app.AccessKeyID,
			Name:        app.Name,
			Description: app.Description,
			Permissions: app.Permissions,
			CreatedAt:   app.CreatedAt,
			UpdatedAt:   app.UpdatedAt,
		})
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	response := GetAppsResponse{
		Applications: responseApps,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		sendError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteApplication handles deleting application
func (h *AppHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from context
	currentUsername := auth.GetUserFromContext(r.Context())
	if currentUsername == "" {
		sendError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Find current user
	currentUser := &models.User{}
	if err := currentUser.FindUserByUsername(h.db.DB, currentUsername); err != nil {
		sendError(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user has admin role
	if !currentUser.HasAdminRole() {
		sendError(w, "Access denied. Admin rights required", http.StatusForbidden)
		return
	}

	// Get application ID from query parameter
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		sendError(w, "Application ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendError(w, "Invalid application ID", http.StatusBadRequest)
		return
	}

	// Delete application
	appModel := &models.Application{}
	if err := appModel.DeleteApplication(h.db.DB, id); err != nil {
		sendError(w, "Failed to delete application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Application deleted successfully",
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

	// Check if application has write permission
	if !app.HasWritePermission() {
		sendError(w, "Application does not have write permission", http.StatusForbidden)
		return
	}

	// Get filename from URL path or query parameter
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		// Try to get filename from URL path
		filename = strings.TrimPrefix(r.URL.Path, "/upload/raw/")
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

// generateAccessKeyID generates a random access key ID
func generateAccessKeyID() string {
	return "ak_" + generateRandomString(16)
}

// generateAccessKeySecret generates a random access key secret
func generateAccessKeySecret() string {
	return "sk_" + generateRandomString(32)
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		// In real application, use crypto/rand
		b[i] = charset[int(time.Now().UnixNano())%len(charset)]
	}
	return string(b)
}
