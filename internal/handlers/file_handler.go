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
)

// ErrorResponse представляет структуру ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// sendError отправляет ошибку в формате JSON
func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message}); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}

// FileRequest представляет структуру JSON запроса для операций с файлами
type FileRequest struct {
	Path string `json:"path"`
}

// CreateDirRequest представляет структуру JSON запроса для создания директории
type CreateDirRequest struct {
	Dirname string `json:"dirname"`
	Path    string `json:"path"`
}

// FileHandler обрабатывает операции с файлами
type FileHandler struct {
	uploadDir string
}

// NewFileHandler создает новый обработчик файлов
func NewFileHandler(uploadDir string) *FileHandler {
	return &FileHandler{
		uploadDir: uploadDir,
	}
}

// UploadFile обрабатывает загрузку файла
func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем Content-Type
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		log.Printf("Invalid Content-Type: %s", contentType)
		sendError(w, "Требуется multipart/form-data", http.StatusBadRequest)
		return
	}

	// Увеличиваем максимальный размер формы и файла
	maxFileSize := int64(32 << 20) // 32MB
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		sendError(w, fmt.Sprintf("Ошибка парсинга формы: %v", err), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		sendError(w, "Ошибка получения файла", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing uploaded file: %v", err)
		}
	}()

	// Проверяем размер файла
	if header.Size > maxFileSize {
		log.Printf("File too large: %d bytes", header.Size)
		sendError(w, fmt.Sprintf("Файл слишком большой. Максимальный размер: %d MB", maxFileSize/(1<<20)), http.StatusBadRequest)
		return
	}

	// Получаем путь для сохранения файла
	path := r.FormValue("path")
	if path == "" {
		path = "."
	}

	// Очищаем и нормализуем путь
	path = filepath.ToSlash(filepath.Clean(path))
	path = regexp.MustCompile(`\s+`).ReplaceAllString(path, "_")
	if strings.Contains(path, "..") {
		log.Printf("Invalid path detected: %s", path)
		sendError(w, "Недопустимый путь", http.StatusBadRequest)
		return
	}

	// Очищаем имя файла от недопустимых символов
	filename := regexp.MustCompile(`[\s\\/:*?"<>|]`).ReplaceAllString(header.Filename, "_")

	// Создаем директорию, если она не существует
	uploadPath := filepath.Join(h.uploadDir, path)
	uploadPath = filepath.Clean(uploadPath)

	// Проверяем, что путь находится внутри разрешенной директории
	uploadDirAbs, _ := filepath.Abs(h.uploadDir)
	fullPath, _ := filepath.Abs(uploadPath)
	if !strings.HasPrefix(fullPath, uploadDirAbs) {
		sendError(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	// Создаем все промежуточные директории
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Printf("Error creating directories: %v", err)
		sendError(w, fmt.Sprintf("Ошибка создания директории: %v", err), http.StatusInternalServerError)
		return
	}

	// Проверяем, что директория действительно создана
	if info, err := os.Stat(uploadPath); err != nil || !info.IsDir() {
		log.Printf("Ошибка проверки директории после создания: %v", err)
		sendError(w, "Ошибка создания директории", http.StatusInternalServerError)
		return
	}

	// Формируем путь для файла с использованием filepath.Join
	filePath := filepath.Join(uploadPath, filename)

	// Проверяем что путь файла находится внутри разрешенной директории
	filePathAbs, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(filePathAbs, uploadDirAbs) {
		sendError(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	// Сохраняем файл
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		sendError(w, "Ошибка создания файла", http.StatusInternalServerError)
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
		sendError(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	log.Printf("Successfully uploaded file %s (%d bytes) to %s", filename, written, path)
	response := map[string]interface{}{
		"message":  "Файл успешно загружен",
		"filename": filename,
		"path":     path,
		"size":     written,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding upload response: %v", err)
	}
}

// UploadRawFile обрабатывает прямую загрузку файла через POST/PUT запрос
func (h *FileHandler) UploadRawFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем имя файла из пути URL
	filename := strings.TrimPrefix(r.URL.Path, "/upload/raw/")
	if filename == "" {
		sendError(w, "Имя файла не указано", http.StatusBadRequest)
		return
	}

	// Создаем файл
	dst, err := os.Create(filepath.Join(h.uploadDir, filename))
	if err != nil {
		sendError(w, "Ошибка создания файла", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			log.Printf("Error closing destination file: %v", err)
		}
	}()

	// Копируем содержимое запроса в файл
	_, err = io.Copy(dst, r.Body)
	if err != nil {
		sendError(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	response := map[string]string{
		"message":  "Файл успешно загружен",
		"filename": filename,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// ListFiles возвращает список файлов и директорий
func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var path string
	if r.Method == http.MethodPost {
		var req FileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding request: %v", err)
			sendError(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}
		path = req.Path
	} else {
		path = r.URL.Query().Get("path")
	}

	if path == "" {
		path = "."
	}

	// Декодируем URL-encoded строку и заменяем недопустимые символы
	path, err := url.QueryUnescape(path)
	if err != nil {
		log.Printf("Error decoding path: %v", err)
		sendError(w, "Некорректный путь", http.StatusBadRequest)
		return
	}

	// Очищаем путь от потенциально опасных символов
	path = filepath.Clean(path)
	if strings.Contains(path, "..") {
		sendError(w, "Недопустимый путь", http.StatusBadRequest)
		return
	}

	// Формируем полный путь
	fullPath := filepath.Join(h.uploadDir, path)

	// Проверяем что путь находится внутри разрешенной директории
	fullPath = filepath.Clean(fullPath)
	uploadDirAbs, _ := filepath.Abs(h.uploadDir)
	fullPath, _ = filepath.Abs(fullPath)

	if !strings.HasPrefix(fullPath, uploadDirAbs) {
		sendError(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	// Проверяем существование директории
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Если директория не найдена, возвращаем пустой массив
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode([]interface{}{}); err != nil {
				log.Printf("Error encoding empty list: %v", err)
			}
			return
		}
		sendError(w, "Ошибка доступа к директории", http.StatusInternalServerError)
		return
	}

	if !info.IsDir() {
		sendError(w, "Указанный путь не является директорией", http.StatusBadRequest)
		return
	}

	files, err := os.ReadDir(fullPath)
	if err != nil {
		// В случае ошибки чтения директории возвращаем пустой массив
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]interface{}{}); err != nil {
			log.Printf("Error encoding empty list: %v", err)
		}
		return
	}

	fileList := make([]map[string]interface{}, 0)

	// Добавляем ссылку на родительскую директорию, если мы не в корне
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
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/download/")
	if filename == "" {
		sendError(w, "Имя файла не указано", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(h.uploadDir, filename)

	// Проверяем существование файла
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		sendError(w, "Файл не найден", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовки для скачивания
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Отдаем файл
	http.ServeFile(w, r, filePath)
}

// CreateDirectory создает новую директорию
func (h *FileHandler) CreateDirectory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req CreateDirRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		sendError(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	// Проверяем наличие имени директории
	if req.Dirname == "" {
		sendError(w, "Необходимо указать имя директории (dirname)", http.StatusBadRequest)
		return
	}

	// Очищаем имя директории от недопустимых символов
	dirName := strings.TrimSpace(req.Dirname)
	// Заменяем недопустимые символы на подчеркивание
	dirName = regexp.MustCompile(`[\s\\/:*?"<>|]`).ReplaceAllString(dirName, "_")

	// Формируем полный путь
	targetPath := dirName
	var pathForJoin string
	if req.Path != "" && req.Path != "." {
		// Очищаем путь и заменяем все обратные слеши на прямые
		pathForJoin = filepath.ToSlash(filepath.Clean(req.Path))
		// Заменяем пробелы в пути на подчеркидения
		pathForJoin = regexp.MustCompile(`\s+`).ReplaceAllString(pathForJoin, "_")
		targetPath = pathForJoin + "/" + dirName
	}

	// Очищаем путь и конвертируем слеши в зависимости от ОС
	targetPath = filepath.FromSlash(targetPath)
	if strings.Contains(targetPath, "..") {
		sendError(w, "Недопустимый путь директории", http.StatusBadRequest)
		return
	}

	// Проверяем что путь находится внутри разрешенной директории
	fullPath := filepath.Join(h.uploadDir, targetPath)
	fullPath = filepath.Clean(fullPath)
	uploadDirAbs, _ := filepath.Abs(h.uploadDir)
	fullPath, _ = filepath.Abs(fullPath)

	if !strings.HasPrefix(fullPath, uploadDirAbs) {
		sendError(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	// Создаем директорию с корректными правами доступа
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		log.Printf("Ошибка создания директории: %v", err)
		sendError(w, fmt.Sprintf("Ошибка создания директории: %v", err), http.StatusInternalServerError)
		return
	}

	// Проверяем что директория действительно создана
	if info, err := os.Stat(fullPath); err != nil || !info.IsDir() {
		log.Printf("Ошибка проверки директории после создания: %v", err)
		sendError(w, "Ошибка создания директории", http.StatusInternalServerError)
		return
	}

	// Получаем относительный путь для ответа
	relativePath, err := filepath.Rel(h.uploadDir, fullPath)
	if err != nil {
		relativePath = targetPath
	}

	// Возвращаем успешный ответ
	response := map[string]string{
		"message": "Директория успешно создана",
		"path":    filepath.ToSlash(relativePath),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding directory creation response: %v", err)
	}
}

// DeleteDirectory удаляет директорию
func (h *FileHandler) DeleteDirectory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		sendError(w, "Путь не указан", http.StatusBadRequest)
		return
	}

	// Очищаем и проверяем путь
	cleanPath := filepath.Clean(path)
	if cleanPath == "." || cleanPath == "/" {
		sendError(w, "Нельзя удалить корневую директорию", http.StatusBadRequest)
		return
	}
	if strings.Contains(cleanPath, "..") {
		sendError(w, "Недопустимый путь", http.StatusBadRequest)
		return
	}

	// Формируем полный путь
	fullPath := filepath.Join(h.uploadDir, cleanPath)

	// Проверяем что путь находится внутри разрешенной директории
	fullPath = filepath.Clean(fullPath)
	uploadDirAbs, _ := filepath.Abs(h.uploadDir)
	fullPath, _ = filepath.Abs(fullPath)
	if !strings.HasPrefix(fullPath, uploadDirAbs) {
		sendError(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	// Проверяем существование директории
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			sendError(w, "Директория не найдена", http.StatusNotFound)
		} else {
			sendError(w, "Ошибка доступа к директории", http.StatusInternalServerError)
		}
		return
	}

	if !info.IsDir() {
		sendError(w, "Указанный путь не является директорией", http.StatusBadRequest)
		return
	}

	// Удаляем директорию
	if err := os.RemoveAll(fullPath); err != nil {
		log.Printf("Error deleting directory: %v", err)
		sendError(w, "Ошибка удаления директории", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	response := map[string]string{
		"message": "Директория успешно удалена",
		"path":    cleanPath,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding directory deletion response: %v", err)
	}
}

// DeleteFile удаляет файл
func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		sendError(w, "Путь не указан", http.StatusBadRequest)
		return
	}

	// Очищаем и проверяем путь
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		sendError(w, "Недопустимый путь", http.StatusBadRequest)
		return
	}

	// Формируем полный путь
	fullPath := filepath.Join(h.uploadDir, cleanPath)

	// Проверяем что путь находится внутри разрешенной директории
	fullPath = filepath.Clean(fullPath)
	uploadDirAbs, _ := filepath.Abs(h.uploadDir)
	fullPath, _ = filepath.Abs(fullPath)
	if !strings.HasPrefix(fullPath, uploadDirAbs) {
		sendError(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	// Проверяем существование файла
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			sendError(w, "Файл не найден", http.StatusNotFound)
		} else {
			sendError(w, "Ошибка доступа к файлу", http.StatusInternalServerError)
		}
		return
	}

	// Проверяем что это файл, а не директория
	if info.IsDir() {
		sendError(w, "Указанный путь является директорией", http.StatusBadRequest)
		return
	}

	// Удаляем файл
	if err := os.Remove(fullPath); err != nil {
		log.Printf("Error deleting file: %v", err)
		sendError(w, "Ошибка удаления файла", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	response := map[string]string{
		"message": "Файл успешно удален",
		"path":    cleanPath,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding file deletion response: %v", err)
	}
}
