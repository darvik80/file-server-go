package main

import (
	"fmt"
	"strings"
)

func main() {
	// Тестирование определения, можно ли отображать файл в браузере
	testContentTypes := []string{
		"image/png",
		"image/jpeg",
		"image/gif",
		"text/html",
		"text/plain",
		"application/pdf",
		"application/json",
		"application/octet-stream",
		"application/zip",
	}
	
	for _, contentType := range testContentTypes {
		isPreviewable := isPreviewableContent(contentType)
		fmt.Printf("Content-Type: %s -> Previewable: %t\n", contentType, isPreviewable)
	}
}

func isPreviewableContent(contentType string) bool {
	isPreviewable := strings.HasPrefix(contentType, "image/") ||
		contentType == "text/html" ||
		contentType == "text/plain" ||
		contentType == "application/pdf" ||
		contentType == "application/json"
	
	return isPreviewable
}