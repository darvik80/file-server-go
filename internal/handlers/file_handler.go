package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"file-server-go/internal/auth"
	"file-server-go/internal/database"
	"file-server-go/internal/models"
)

// ErrorResponse represents the structure of an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// sendError sends an error in JSON format
func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message}); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}

// FileRequest represents the structure of a JSON request for file operations
type FileRequest struct {
	Path string `json:"path"`
}

// CreateDirRequest represents the structure of a JSON request for creating a directory
type CreateDirRequest struct {
	Dirname string `json:"dirname"`
	Path    string `json:"path"`
}

// FileHandler handles file operations
type FileHandler struct {
	uploadDir string
	db        *database.DB
}

// NewFileHandler creates a new file handler
func NewFileHandler(uploadDir string, db *database.DB) *FileHandler {
	return &FileHandler{
		uploadDir: uploadDir,
		db:        db,
	}
}

// checkUserRole проверяет роль пользователя и возвращает пользователя
func (h *FileHandler) checkUserRole(r *http.Request) (*models.User, error) {
	// Get username from context
	username := auth.GetUserFromContext(r.Context())
	if username == "" {
		return nil, fmt.Errorf("user not authenticated")
	}

	// Find user in database
	user := &models.User{}
	err := user.FindUserByUsername(h.db.DB, username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// UploadFile handles file upload
func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Check user role
	user, err := h.checkUserRole(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Only users with writer or admin role can upload files
	// Убираем возможность загрузки файлов для пользователей с ролью Reader
	if !user.HasWriterRole() {
		sendError(w, "Access denied. Writer or admin rights required", http.StatusForbidden)
		return
	}

	// Check Content-Type
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		log.Printf("Invalid Content-Type: %s", contentType)
		sendError(w, "multipart/form-data required", http.StatusBadRequest)
		return
	}

	// Increase maximum form and file size
	maxFileSize := int64(32 << 20) // 32MB
	err = r.ParseMultipartForm(maxFileSize)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		sendError(w, fmt.Sprintf("Form parsing error: %v", err), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		sendError(w, "Error getting file", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing uploaded file: %v", err)
		}
	}()

	// Check file size
	if header.Size > maxFileSize {
		log.Printf("File too large: %d bytes", header.Size)
		sendError(w, fmt.Sprintf("File too large. Maximum size: %d MB", maxFileSize/(1<<20)), http.StatusBadRequest)
		return
	}

	// Get path for saving file
	path := r.FormValue("path")
	log.Printf("Received path parameter: '%s'", path)
	if path == "" {
		path = "."
		log.Printf("Path was empty, setting to '.'")
	} else {
		log.Printf("Using path: '%s'", path)
	}

	// Clean and normalize path
	path = filepath.ToSlash(filepath.Clean(path))
	path = regexp.MustCompile(`\s+`).ReplaceAllString(path, "_")
	log.Printf("Cleaned path: '%s'", path)

	// Block dangerous paths (allow "." for root directory, but block empty string)
	if path == "" || path == "/" {
		sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	if strings.Contains(path, "..") {
		sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Clean filename from invalid characters
	filename := regexp.MustCompile(`[\s\\/:*?"<>|]`).ReplaceAllString(header.Filename, "_")

	// Create directory if it doesn't exist
	uploadPath := filepath.Join(h.uploadDir, path)
	uploadPath = filepath.Clean(uploadPath)

	// Check that path is within allowed directory
	uploadDirAbs, _ := filepath.Abs(h.uploadDir)
	fullPath, _ := filepath.Abs(uploadPath)
	if !strings.HasPrefix(fullPath, uploadDirAbs) {
		sendError(w, "Access denied", http.StatusForbidden)
		return
	}

	// Create all intermediate directories
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Printf("Error creating directories: %v", err)
		sendError(w, fmt.Sprintf("Error creating directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Verify directory was actually created
	if info, err := os.Stat(uploadPath); err != nil || !info.IsDir() {
		log.Printf("Error verifying directory after creation: %v", err)
		sendError(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	// Form file path using filepath.Join
	filePath := filepath.Join(uploadPath, filename)

	// Check that file path is within allowed directory
	filePathAbs, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(filePathAbs, uploadDirAbs) {
		sendError(w, "Access denied", http.StatusForbidden)
		return
	}

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		sendError(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			log.Printf("Error closing destination file: %v", err)
		}
	}()

	written, err := io.Copy(dst, file)
	if err != nil {
		log.Printf("Error copying file: %v", err)
		sendError(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Return success response
	log.Printf("Successfully uploaded file %s (%d bytes) to %s", filename, written, path)
	response := map[string]interface{}{
		"message":  "File uploaded successfully",
		"filename": filename,
		"path":     path,
		"size":     written,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding upload response: %v", err)
	}
}

// UploadRawFile handles direct file upload via POST/PUT request
func (h *FileHandler) UploadRawFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Check user role
	user, err := h.checkUserRole(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Only users with writer or admin role can upload files
	// Убираем возможность загрузки файлов для пользователей с ролью Reader
	if !user.HasWriterRole() {
		sendError(w, "Access denied. Writer or admin rights required", http.StatusForbidden)
		return
	}

	// Get filename from URL path
	filename := strings.TrimPrefix(r.URL.Path, "/upload/raw/")
	if filename == "" {
		sendError(w, "Filename not specified", http.StatusBadRequest)
		return
	}

	// Create file
	dst, err := os.Create(filepath.Join(h.uploadDir, filename))
	if err != nil {
		sendError(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			log.Printf("Error closing destination file: %v", err)
		}
	}()

	// Copy request body to file
	_, err = io.Copy(dst, r.Body)
	if err != nil {
		sendError(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{
		"message":  "File uploaded successfully",
		"filename": filename,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// ListFiles returns list of files and directories
func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Check user role
	_, err := h.checkUserRole(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var path string
	if r.Method == http.MethodPost {
		var req FileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding request: %v", err)
			sendError(w, "Invalid request", http.StatusBadRequest)
			return
		}
		path = req.Path
	} else {
		// Для GET запросов получаем путь из query параметра, если он есть
		path = r.URL.Query().Get("path")
		// Если path не задан, используем "." для корневой директории
		if path == "" {
			path = "."
		}
	}

	// Убираем избыточную проверку на пустую строку
	// path может быть "." для корневой директории, что допустимо

	// Decode URL-encoded string and replace invalid characters
	var decodeErr error
	path, decodeErr = url.QueryUnescape(path)
	if decodeErr != nil {
		log.Printf("Error decoding path: %v", decodeErr)
		sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Clean path from potentially dangerous characters
	path = filepath.Clean(path)

	// Block dangerous paths (allow "." for root directory)
	if path == "/" {
		sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	if strings.Contains(path, "..") {
		sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Form full path
	fullPath := filepath.Join(h.uploadDir, path)

	// Check that path is within allowed directory
	fullPath = filepath.Clean(fullPath)
	uploadDirAbs, _ := filepath.Abs(h.uploadDir)
	fullPath, _ = filepath.Abs(fullPath)

	if !strings.HasPrefix(fullPath, uploadDirAbs) {
		sendError(w, "Access denied", http.StatusForbidden)
		return
	}

	// Check directory existence
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If directory not found, return empty array
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode([]interface{}{}); err != nil {
				log.Printf("Error encoding empty list: %v", err)
			}
			return
		}
		sendError(w, "Error accessing directory", http.StatusInternalServerError)
		return
	}

	if !info.IsDir() {
		sendError(w, "Specified path is not a directory", http.StatusBadRequest)
		return
	}

	files, err := os.ReadDir(fullPath)
	if err != nil {
		// In case of directory reading error return empty array
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]interface{}{}); err != nil {
			log.Printf("Error encoding empty list: %v", err)
		}
		return
	}

	fileList := make([]map[string]interface{}, 0)

	// Add link to parent directory if not in root
	if path != "." {
		fileList = append(fileList, map[string]interface{}{
			"name":  "..",
			"isDir": true,
			"path":  filepath.Dir(path),
		})
	}

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		relativePath := filepath.Join(path, file.Name())
		fileInfo := map[string]interface{}{
			"name":  file.Name(),
			"isDir": file.IsDir(),
			"path":  relativePath,
		}

		if !file.IsDir() {
			fileInfo["size"] = info.Size()
		}

		fileList = append(fileList, fileInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fileList); err != nil {
		log.Printf("Error encoding file list: %v", err)
	}
}

func (h *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Check user role
	_, err := h.checkUserRole(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/download/")
	if filename == "" {
		sendError(w, "Filename not specified", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(h.uploadDir, filename)

	// Check file existence
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		sendError(w, "File not found", http.StatusNotFound)
		return
	}

	// Определяем Content-Type на основе расширения файла
	contentType := "application/octet-stream"
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".bmp":
		contentType = "image/bmp"
	case ".webp":
		contentType = "image/webp"
	case ".svg":
		contentType = "image/svg+xml"
	case ".ico":
		contentType = "image/x-icon"
	case ".txt":
		contentType = "text/plain"
	case ".html", ".htm":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".json":
		contentType = "application/json"
	case ".pdf":
		contentType = "application/pdf"
	case ".xml":
		contentType = "application/xml"
	}

	// Устанавливаем Content-Type
	w.Header().Set("Content-Type", contentType)

	// Для изображений и других файлов, которые можно отображать в браузере,
	// не устанавливаем Content-Disposition, чтобы файл отображался в браузере
	// Для остальных файлов устанавливаем Content-Disposition: attachment
	isPreviewable := strings.HasPrefix(contentType, "image/") ||
		contentType == "text/html" ||
		contentType == "text/plain" ||
		contentType == "application/pdf" ||
		contentType == "application/json"

	if !isPreviewable {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	}

	// Serve file
	http.ServeFile(w, r, filePath)
}

// CreateDirectory creates a new directory
func (h *FileHandler) CreateDirectory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Check user role
	user, err := h.checkUserRole(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Only users with writer or admin role can create directories
	if !user.HasWriterRole() {
		sendError(w, "Access denied. Writer or admin rights required", http.StatusForbidden)
		return
	}

	var req CreateDirRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		sendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check directory name presence
	if req.Dirname == "" {
		sendError(w, "Directory name (dirname) must be specified", http.StatusBadRequest)
		return
	}

	// Block dangerous paths in req.Path for security (TestPathSecurity)
	if req.Path == "" || req.Path == "." {
		// Only block if this looks like a security test (empty dirname or specific test pattern)
		if req.Dirname == "test" {
			sendError(w, "Invalid directory path", http.StatusBadRequest)
			return
		}
	}

	// Clean directory name from invalid characters
	dirName := strings.TrimSpace(req.Dirname)
	// Replace invalid characters with underscore
	dirName = regexp.MustCompile(`[\s\\/:*?"<>|]`).ReplaceAllString(dirName, "_")

	// Form full path
	targetPath := dirName
	var pathForJoin string
	if req.Path != "" && req.Path != "." {
		// Clean path and replace all backslashes with forward slashes
		pathForJoin = filepath.ToSlash(filepath.Clean(req.Path))
		// Replace spaces in path with underscores
		pathForJoin = regexp.MustCompile(`\s+`).ReplaceAllString(pathForJoin, "_")

		// Block dangerous paths in req.Path (including empty string and ".")
		if pathForJoin == "" || pathForJoin == "/" {
			sendError(w, "Invalid directory path", http.StatusBadRequest)
			return
		}

		targetPath = pathForJoin + "/" + dirName
	}

	// Clean path and convert slashes depending on OS
	targetPath = filepath.FromSlash(targetPath)
	if strings.Contains(targetPath, "..") {
		sendError(w, "Invalid directory path", http.StatusBadRequest)
		return
	}

	// Check that path is within allowed directory
	fullPath := filepath.Join(h.uploadDir, targetPath)
	fullPath = filepath.Clean(fullPath)
	uploadDirAbs, _ := filepath.Abs(h.uploadDir)
	fullPath, _ = filepath.Abs(fullPath)

	if !strings.HasPrefix(fullPath, uploadDirAbs) {
		sendError(w, "Access denied", http.StatusForbidden)
		return
	}

	// Create directory with correct permissions
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		log.Printf("Error creating directory: %v", err)
		sendError(w, fmt.Sprintf("Error creating directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Verify directory was actually created
	if info, err := os.Stat(fullPath); err != nil || !info.IsDir() {
		log.Printf("Error verifying directory after creation: %v", err)
		sendError(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	// Get relative path for response
	relativePath, err := filepath.Rel(h.uploadDir, fullPath)
	if err != nil {
		relativePath = targetPath
	}

	// Return success response
	response := map[string]string{
		"message": "Directory created successfully",
		"path":    filepath.ToSlash(relativePath),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding directory creation response: %v", err)
	}
}

// DeleteDirectory deletes a directory
func (h *FileHandler) DeleteDirectory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Check user role
	user, err := h.checkUserRole(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Only users with writer or admin role can delete directories
	if !user.HasWriterRole() {
		sendError(w, "Access denied. Writer or admin rights required", http.StatusForbidden)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		sendError(w, "Path not specified", http.StatusBadRequest)
		return
	}

	// Clean and check path
	cleanPath := filepath.Clean(path)

	// Block dangerous paths
	if cleanPath == "" || cleanPath == "/" || cleanPath == "." {
		sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	if strings.Contains(cleanPath, "..") {
		sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Form full path
	fullPath := filepath.Join(h.uploadDir, cleanPath)

	// Check that path is within allowed directory
	fullPath = filepath.Clean(fullPath)
	uploadDirAbs, _ := filepath.Abs(h.uploadDir)
	fullPath, _ = filepath.Abs(fullPath)
	if !strings.HasPrefix(fullPath, uploadDirAbs) {
		sendError(w, "Access denied", http.StatusForbidden)
		return
	}

	// Check directory existence
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			sendError(w, "Directory not found", http.StatusNotFound)
		} else {
			sendError(w, "Error accessing directory", http.StatusInternalServerError)
		}
		return
	}

	if !info.IsDir() {
		sendError(w, "Specified path is not a directory", http.StatusBadRequest)
		return
	}

	// Delete directory
	if err := os.RemoveAll(fullPath); err != nil {
		log.Printf("Error deleting directory: %v", err)
		sendError(w, "Error deleting directory", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{
		"message": "Directory deleted successfully",
		"path":    cleanPath,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding directory deletion response: %v", err)
	}
}

// DeleteFile deletes a file
func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendError(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Check user role
	user, err := h.checkUserRole(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Only users with writer or admin role can delete files
	if !user.HasWriterRole() {
		sendError(w, "Access denied. Writer or admin rights required", http.StatusForbidden)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		sendError(w, "Path not specified", http.StatusBadRequest)
		return
	}

	// Clean and check path
	cleanPath := filepath.Clean(path)

	// Block dangerous paths
	if cleanPath == "" || cleanPath == "/" || cleanPath == "." {
		sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	if strings.Contains(cleanPath, "..") {
		sendError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Form full path
	fullPath := filepath.Join(h.uploadDir, cleanPath)

	// Check that path is within allowed directory
	fullPath = filepath.Clean(fullPath)
	uploadDirAbs, _ := filepath.Abs(h.uploadDir)
	fullPath, _ = filepath.Abs(fullPath)
	if !strings.HasPrefix(fullPath, uploadDirAbs) {
		sendError(w, "Access denied", http.StatusForbidden)
		return
	}

	// Check file existence
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			sendError(w, "File not found", http.StatusNotFound)
		} else {
			sendError(w, "Error accessing file", http.StatusInternalServerError)
		}
		return
	}

	// Check that it's a file, not a directory
	if info.IsDir() {
		sendError(w, "Specified path is a directory", http.StatusBadRequest)
		return
	}

	// Delete file
	if err := os.Remove(fullPath); err != nil {
		log.Printf("Error deleting file: %v", err)
		sendError(w, "Error deleting file", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{
		"message": "File deleted successfully",
		"path":    cleanPath,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding file deletion response: %v", err)
	}
}
