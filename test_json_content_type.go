package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	// Тестирование определения типа контента для JSON файла
	filename := "jsconfig.json"
	
	contentType := getContentType(filename)
	isPreviewable := isPreviewableContent(contentType)
	
	fmt.Printf("Файл: %s\n", filename)
	fmt.Printf("Content-Type: %s\n", contentType)
	fmt.Printf("Previewable: %t\n", isPreviewable)
}

func getContentType(filename string) string {
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
	
	return contentType
}

func isPreviewableContent(contentType string) bool {
	isPreviewable := strings.HasPrefix(contentType, "image/") ||
		contentType == "text/html" ||
		contentType == "text/plain" ||
		contentType == "application/pdf" ||
		contentType == "application/json"
	
	return isPreviewable
}