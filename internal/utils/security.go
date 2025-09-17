package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidatePath validates and cleans file path for security
func ValidatePath(path, uploadDir string) (string, error) {
	// Clean path
	cleanPath := filepath.Clean(path)
	
	// Block dangerous paths
	if cleanPath == "/" || strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("invalid path")
	}
	
	// Form full path
	fullPath := filepath.Join(uploadDir, cleanPath)
	fullPath = filepath.Clean(fullPath)
	
	// Check that path is within allowed directory
	uploadDirAbs, _ := filepath.Abs(uploadDir)
	fullPathAbs, _ := filepath.Abs(fullPath)
	
	if !strings.HasPrefix(fullPathAbs, uploadDirAbs) {
		return "", fmt.Errorf("access denied")
	}
	
	return fullPath, nil
}

// IsPreviewableContent determines if content can be previewed in browser
func IsPreviewableContent(contentType string) bool {
	return strings.HasPrefix(contentType, "image/") ||
		contentType == "text/html" ||
		contentType == "text/plain" ||
		contentType == "application/pdf" ||
		contentType == "application/json"
}

// GetContentTypeByExtension returns content type based on file extension
func GetContentTypeByExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".txt":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".pdf":
		return "application/pdf"
	case ".xml":
		return "application/xml"
	default:
		return "application/octet-stream"
	}
}