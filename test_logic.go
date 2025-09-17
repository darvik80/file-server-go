package main

import (
	"fmt"
	"strings"
)

func main() {
	// Тестирование логики определения поддерживаемых типов для предварительного просмотра
	contentType := "application/json"
	
	isPreviewable := strings.HasPrefix(contentType, "image/") ||
		contentType == "text/html" ||
		contentType == "text/plain" ||
		contentType == "application/pdf" ||
		contentType == "application/json"
	
	fmt.Printf("Content-Type: %s\n", contentType)
	fmt.Printf("Previewable: %t\n", isPreviewable)
	
	if !isPreviewable {
		fmt.Println("Content-Disposition would be set to attachment")
	} else {
		fmt.Println("Content-Disposition would NOT be set")
	}
}