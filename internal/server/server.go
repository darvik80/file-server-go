package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"file-server-go/internal/auth"
	"file-server-go/internal/config"
	"file-server-go/internal/handlers"
	"file-server-go/internal/middleware"
)

// Server представляет HTTP сервер
type Server struct {
	config      *config.Config
	fileHandler *handlers.FileHandler
	authHandler *handlers.AuthHandler
	templates   *template.Template
}

// New создает новый сервер
func New(cfg *config.Config) *Server {
	// Устанавливаем ключ для JWT
	auth.SetJWTKey(cfg.JWTKey)

	templates := template.Must(template.ParseFiles("web/templates/index.html"))

	return &Server{
		config:      cfg,
		fileHandler: handlers.NewFileHandler(cfg.UploadDir),
		authHandler: handlers.NewAuthHandler(cfg),
		templates:   templates,
	}
}

// Router настраивает и возвращает HTTP роутер
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	// Маршруты авторизации
	mux.HandleFunc("/login", s.authHandler.Login)

	// Маршруты для работы с файлами
	mux.HandleFunc("/upload", s.fileHandler.UploadFile)
	mux.HandleFunc("/upload/raw/", s.fileHandler.UploadRawFile)
	mux.HandleFunc("/files", s.fileHandler.ListFiles)
	mux.HandleFunc("/download/", s.fileHandler.DownloadFile)
	mux.HandleFunc("/create-dir", s.fileHandler.CreateDirectory)
	mux.HandleFunc("/delete-dir", s.fileHandler.DeleteDirectory)
	mux.HandleFunc("/delete-file", s.fileHandler.DeleteFile)

	// Статические файлы
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js/"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css/"))))

	// Главная страница
	mux.HandleFunc("/", s.handleHome)

	// Оборачиваем обработчики в middleware в правильном порядке
	var handler http.Handler = mux
	handler = handlePanic(handler)        // Первым идет обработка паники
	handler = middleware.Logging(handler) // Затем логирование
	handler = middleware.Auth(handler)    // Потом авторизация
	handler = middleware.CORS(handler)    // И последним CORS

	return handler
}

// handlePanic восстанавливает работу после паники
func handlePanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Internal Server Error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// handleHome обрабатывает главную страницу
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
