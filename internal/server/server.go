package server

import (
	"html/template"
	"net/http"

	"file-server-go/internal/assets"
	"file-server-go/internal/auth"
	"file-server-go/internal/config"
	"file-server-go/internal/database"
	"file-server-go/internal/handlers"
	"file-server-go/internal/middleware"
)

// Server represents HTTP server
type Server struct {
	config      *config.Config
	fileHandler *handlers.FileHandler
	authHandler *handlers.AuthHandler
	appHandler  *handlers.AppHandler
	userHandler *handlers.UserHandler
	templates   *template.Template
	db          *database.DB
}

// New creates new server
func New(cfg *config.Config, fileHandler *handlers.FileHandler, db *database.DB) *Server {
	// Set JWT key
	auth.SetJWTKey(cfg.JWTKey)

	// Load templates from embedded files
	templateFS := assets.GetTemplateFS()
	templates, err := template.ParseFS(templateFS, "*.html")
	if err != nil {
		panic("Error loading templates: " + err.Error())
	}

	return &Server{
		config:      cfg,
		fileHandler: handlers.NewFileHandler(cfg.UploadDir, db),
		authHandler: handlers.NewAuthHandler(cfg, db),
		appHandler:  handlers.NewAppHandler(db, cfg.UploadDir),
		userHandler: handlers.NewUserHandler(db),
		templates:   templates,
		db:          db,
	}
}

// Router configures and returns HTTP router
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	// Authentication route
	mux.HandleFunc("/login", s.authHandler.Login)
	mux.HandleFunc("/user-info", s.authHandler.GetUserInfo)

	// Application routes
	mux.HandleFunc("/register-app", s.appHandler.RegisterApp)
	mux.HandleFunc("/applications", s.appHandler.GetApplications)
	mux.HandleFunc("/delete-app", s.appHandler.DeleteApplication)

	// User routes
	mux.HandleFunc("/users", s.userHandler.GetUsers)
	mux.HandleFunc("/create-user", s.userHandler.CreateUser)
	mux.HandleFunc("/delete-user", s.userHandler.DeleteUser)

	// Application raw file upload route
	mux.HandleFunc("/upload/raw", s.appHandler.UploadRawFile)
	mux.HandleFunc("/upload/raw/", s.appHandler.UploadRawFile)

	// File operations routes
	mux.HandleFunc("/upload", s.fileHandler.UploadFile)
	mux.HandleFunc("/files", s.fileHandler.ListFiles)
	mux.HandleFunc("/download/", s.fileHandler.DownloadFile)
	mux.HandleFunc("/create-dir", s.fileHandler.CreateDirectory)
	mux.HandleFunc("/delete-dir", s.fileHandler.DeleteDirectory)
	mux.HandleFunc("/delete-file", s.fileHandler.DeleteFile)

	// Embedded static files
	staticFS := assets.GetStaticFS()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Main page using embedded template
	mux.HandleFunc("/", s.handleHome)

	// Apply middleware in correct order
	handler := middleware.CORS(mux)
	handler = middleware.Logging(handler)
	handler = middleware.AppAuth(s.db)(handler) // Application authentication middleware
	handler = middleware.Auth(handler)          // User authentication middleware

	return handler
}

// handleHome handles main page using embedded template
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Use embedded template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := struct {
		Title string
	}{
		Title: "File Server",
	}

	if err := s.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
		return
	}
}
