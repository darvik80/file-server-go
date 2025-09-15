package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileHandler(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "file_handler_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	handler := NewFileHandler(tempDir)

	t.Run("UploadFile", func(t *testing.T) {
		testUploadFile(t, handler, tempDir)
	})

	t.Run("ListFiles", func(t *testing.T) {
		testListFiles(t, handler, tempDir)
	})

	t.Run("DownloadFile", func(t *testing.T) {
		testDownloadFile(t, handler, tempDir)
	})

	t.Run("CreateDirectory", func(t *testing.T) {
		testCreateDirectory(t, handler, tempDir)
	})

	t.Run("DeleteDirectory", func(t *testing.T) {
		testDeleteDirectory(t, handler, tempDir)
	})

	t.Run("DeleteFile", func(t *testing.T) {
		testDeleteFile(t, handler, tempDir)
	})

	t.Run("PreviewHTMLFile", func(t *testing.T) {
		testPreviewHTMLFile(t, handler, tempDir)
	})

	t.Run("ErrorCases", func(t *testing.T) {
		testErrorCases(t, handler, tempDir)
	})
}

func testUploadFile(t *testing.T, handler *FileHandler, tempDir string) {
	tests := []struct {
		name          string
		filename      string
		path          string
		content       string
		expectedError bool
	}{
		{
			name:     "valid file upload",
			filename: "test.txt",
			path:     ".",
			content:  "Hello World",
		},
		{
			name:     "upload with subdirectory",
			filename: "test2.txt",
			path:     "subdir",
			content:  "Test content",
		},
		{
			name:          "invalid path",
			filename:      "test.txt",
			path:          "../invalid",
			content:       "Test",
			expectedError: true,
		},
		{
			name:     "html file upload",
			filename: "test.html",
			path:     ".",
			content:  "<html><body><h1>Test HTML</h1></body></html>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// Добавляем поле path
			if err := writer.WriteField("path", tt.path); err != nil {
				t.Fatalf("Failed to write path field: %v", err)
			}

			// Добавляем файл
			part, err := writer.CreateFormFile("file", tt.filename)
			if err != nil {
				t.Fatalf("Failed to create form file: %v", err)
			}
			part.Write([]byte(tt.content))

			writer.Close()

			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

			handler.UploadFile(w, req)

			resp := w.Result()
			if tt.expectedError {
				if resp.StatusCode == http.StatusOK {
					t.Errorf("Expected error but got success for test: %s", tt.name)
				}
				return
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d for test: %s", resp.StatusCode, tt.name)
			}

			// Проверяем что файл создан
			fullPath := filepath.Join(tempDir, tt.path, tt.filename)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Errorf("File was not created: %s", fullPath)
			}
		})
	}
}

func testListFiles(t *testing.T, handler *FileHandler, tempDir string) {
	// Создаем тестовые файлы и директории
	testFiles := []struct {
		path    string
		isDir   bool
		content string
	}{
		{"file1.txt", false, "content1"},
		{"subdir", true, ""},
		{"subdir/file2.txt", false, "content2"},
		{"test.html", false, "<html><body><h1>Test</h1></body></html>"},
	}

	for _, tf := range testFiles {
		fullPath := filepath.Join(tempDir, tf.path)
		if tf.isDir {
			os.MkdirAll(fullPath, 0755)
		} else {
			os.MkdirAll(filepath.Dir(fullPath), 0755)
			os.WriteFile(fullPath, []byte(tf.content), 0644)
		}
	}

	tests := []struct {
		name     string
		path     string
		method   string
		expected int // ожидаемое количество файлов
	}{
		{"root directory", ".", "GET", 4},    // file1.txt + subdir + test.html
		{"subdirectory", "subdir", "GET", 3}, // file2.txt + ..
		{"post method", ".", "POST", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.method == "POST" {
				body, _ := json.Marshal(FileRequest{Path: tt.path})
				req = httptest.NewRequest("POST", "/files", bytes.NewReader(body))
			} else {
				req = httptest.NewRequest("GET", fmt.Sprintf("/files?path=%s", url.QueryEscape(tt.path)), nil)
			}

			w := httptest.NewRecorder()
			handler.ListFiles(w, req)

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}

			var files []map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if len(files) != tt.expected {
				t.Errorf("Expected %d files, got %d", tt.expected, len(files))
			}
		})
	}
}

func testDownloadFile(t *testing.T, handler *FileHandler, tempDir string) {
	// Создаем тестовый файл
	testContent := "Hello, this is test content for download"
	testFile := "download_test.txt"
	fullPath := filepath.Join(tempDir, testFile)
	os.WriteFile(fullPath, []byte(testContent), 0644)

	t.Run("valid download", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/download/"+testFile, nil)
		w := httptest.NewRecorder()

		handler.DownloadFile(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		content, _ := io.ReadAll(resp.Body)
		if string(content) != testContent {
			t.Errorf("Downloaded content doesn't match expected")
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/download/nonexistent.txt", nil)
		w := httptest.NewRecorder()

		handler.DownloadFile(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func testCreateDirectory(t *testing.T, handler *FileHandler, tempDir string) {
	tests := []struct {
		name          string
		dirname       string
		path          string
		expectedError bool
	}{
		{"create in root", "newdir1", ".", false},
		{"create in subdir", "newdir2", "subdir", false},
		{"invalid chars in name", "invalid/name", ".", false},
		{"empty dirname", "", ".", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(CreateDirRequest{
				Dirname: tt.dirname,
				Path:    tt.path,
			})

			req := httptest.NewRequest("POST", "/create-dir", bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateDirectory(w, req)

			resp := w.Result()
			if tt.expectedError {
				if resp.StatusCode == http.StatusOK {
					t.Errorf("Expected error but got success for test: %s", tt.name)
				}
				return
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d for test: %s", resp.StatusCode, tt.name)
			}

			// Проверяем что директория создана
			tmpDir := strings.ReplaceAll(tt.dirname, "/", "_")
			fullPath := filepath.Join(tempDir, tt.path, tmpDir)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Errorf("Directory was not created: %s", fullPath)
			}
		})
	}
}

func testDeleteDirectory(t *testing.T, handler *FileHandler, tempDir string) {
	// Создаем тестовую директорию
	testDir := "dir_to_delete"
	fullPath := filepath.Join(tempDir, testDir)
	os.MkdirAll(fullPath, 0755)
	os.WriteFile(filepath.Join(fullPath, "test.txt"), []byte("content"), 0644)

	t.Run("valid delete", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/delete-dir?path="+testDir, nil)
		w := httptest.NewRecorder()

		handler.DeleteDirectory(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Проверяем что директория удалена
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			t.Errorf("Directory was not deleted: %s", fullPath)
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/delete-dir?path=nonexistent", nil)
		w := httptest.NewRecorder()

		handler.DeleteDirectory(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("delete root directory", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/delete-dir?path=.", nil)
		w := httptest.NewRecorder()

		handler.DeleteDirectory(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func testDeleteFile(t *testing.T, handler *FileHandler, tempDir string) {
	// Создаем тестовый файл
	testFile := "file_to_delete.txt"
	fullPath := filepath.Join(tempDir, testFile)
	os.WriteFile(fullPath, []byte("content"), 0644)

	t.Run("valid delete", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/delete-file?path="+testFile, nil)
		w := httptest.NewRecorder()

		handler.DeleteFile(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Проверяем что файл удален
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			t.Errorf("File was not deleted: %s", fullPath)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/delete-file?path=nonexistent.txt", nil)
		w := httptest.NewRecorder()

		handler.DeleteFile(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("delete directory as file", func(t *testing.T) {
		// Создаем тестовую директорию
		testDir := "dir_as_file"
		fullPath := filepath.Join(tempDir, testDir)
		os.MkdirAll(fullPath, 0755)

		req := httptest.NewRequest("DELETE", "/delete-file?path="+testDir, nil)
		w := httptest.NewRecorder()

		handler.DeleteFile(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

func testPreviewHTMLFile(t *testing.T, handler *FileHandler, tempDir string) {
	// Создаем тестовый HTML файл
	htmlContent := "<html><body><h1>Test HTML Preview</h1></body></html>"
	htmlFile := "test.html"
	fullPath := filepath.Join(tempDir, htmlFile)
	os.WriteFile(fullPath, []byte(htmlContent), 0644)

	t.Run("valid html preview", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/download/"+htmlFile, nil)
		w := httptest.NewRecorder()

		handler.DownloadFile(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		content, _ := io.ReadAll(resp.Body)
		if string(content) != htmlContent {
			t.Errorf("HTML content doesn't match expected")
		}

		// Проверяем заголовки
		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/octet-stream" {
			t.Errorf("Expected Content-Type application/octet-stream, got %s", contentType)
		}

		contentDisposition := resp.Header.Get("Content-Disposition")
		if !strings.Contains(contentDisposition, "attachment; filename=") {
			t.Errorf("Expected Content-Disposition with attachment, got %s", contentDisposition)
		}
	})

	t.Run("non-existent html file", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/download/nonexistent.html", nil)
		w := httptest.NewRecorder()

		handler.DownloadFile(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func testErrorCases(t *testing.T, handler *FileHandler, tempDir string) {
	tests := []struct {
		name     string
		method   string
		endpoint string
		body     io.Reader
		expected int
	}{
		{"wrong method upload", "GET", "/upload", nil, http.StatusMethodNotAllowed},
		{"wrong method list", "PUT", "/files", nil, http.StatusMethodNotAllowed},
		{"wrong method download", "POST", "/download/test", nil, http.StatusMethodNotAllowed},
		{"wrong method create dir", "GET", "/create-dir", nil, http.StatusMethodNotAllowed},
		{"wrong method delete dir", "POST", "/delete-dir?path=test", nil, http.StatusMethodNotAllowed},
		{"wrong method delete file", "POST", "/delete-file?path=test", nil, http.StatusMethodNotAllowed},
		{"invalid json create dir", "POST", "/create-dir", strings.NewReader("{invalid}"), http.StatusBadRequest},
		{"missing file in upload", "POST", "/upload", strings.NewReader(""), http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.endpoint, tt.body)
			if tt.body != nil && tt.method == "POST" {
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()

			switch tt.endpoint {
			case "/upload":
				handler.UploadFile(w, req)
			case "/files":
				handler.ListFiles(w, req)
			case "/download/test":
				handler.DownloadFile(w, req)
			case "/create-dir":
				handler.CreateDirectory(w, req)
			case "/delete-dir?path=test":
				handler.DeleteDirectory(w, req)
			case "/delete-file?path=test":
				handler.DeleteFile(w, req)
			}

			resp := w.Result()
			if resp.StatusCode != tt.expected {
				t.Errorf("Expected status %d, got %d for test: %s", tt.expected, resp.StatusCode, tt.name)
			}
		})
	}
}

func TestPathSecurity(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "security_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	handler := NewFileHandler(tempDir)

	// Тестируем различные попытки обхода безопасности
	maliciousPaths := []string{
		"../outside",
		"../../etc/passwd",
		"..\\windows\\system32",
		"subdir/../../root",
		"subdir/../..",
		"",
		"/",
		".",
	}

	for _, path := range maliciousPaths {
		t.Run("path_traversal_"+path, func(t *testing.T) {
			// Тестируем ListFiles (разрешаем "." как валидный путь для корневой директории)
			if path != "." {
				req := httptest.NewRequest("GET", "/files?path="+url.QueryEscape(path), nil)
				w := httptest.NewRecorder()
				handler.ListFiles(w, req)

				resp := w.Result()
				if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusForbidden {
					t.Errorf("Expected 400/403 for path %s, got %d", path, resp.StatusCode)
				}
			}

			// Тестируем CreateDirectory
			body, _ := json.Marshal(CreateDirRequest{Dirname: "test", Path: path})
			req := httptest.NewRequest("POST", "/create-dir", bytes.NewReader(body))
			w := httptest.NewRecorder()
			handler.CreateDirectory(w, req)

			resp := w.Result()
			if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusForbidden {
				t.Errorf("Expected 400/403 for create dir with path %s, got %d", path, resp.StatusCode)
			}

			// Тестируем DeleteDirectory
			req = httptest.NewRequest("DELETE", "/delete-dir?path="+path, nil)
			w = httptest.NewRecorder()
			handler.DeleteDirectory(w, req)

			resp = w.Result()
			if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusForbidden {
				t.Errorf("Expected 400/403 for delete dir with path %s, got %d", path, resp.StatusCode)
			}

			// Тестируем DeleteFile
			req = httptest.NewRequest("DELETE", "/delete-file?path="+path, nil)
			w = httptest.NewRecorder()
			handler.DeleteFile(w, req)

			resp = w.Result()
			if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusForbidden {
				t.Errorf("Expected 400/403 for delete file with path %s, got %d", path, resp.StatusCode)
			}
		})
	}
}

func TestSendError(t *testing.T) {
	w := httptest.NewRecorder()
	sendError(w, "Test error", http.StatusBadRequest)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	var errorResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errorResp.Error != "Test error" {
		t.Errorf("Expected error message 'Test error', got '%s'", errorResp.Error)
	}
}

// Benchmark тесты для измерения производительности
func BenchmarkUploadFile(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "benchmark_test")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	handler := NewFileHandler(tempDir)

	for i := 0; i < b.N; i++ {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("path", ".")

		part, _ := writer.CreateFormFile("file", fmt.Sprintf("test%d.txt", i))
		part.Write([]byte("benchmark content"))
		writer.Close()

		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		handler.UploadFile(w, req)
	}
}
