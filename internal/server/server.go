package server

import (
	"net/http"

	"file-server-go/internal/config"
	"file-server-go/internal/handlers"
	"file-server-go/internal/middleware"
)

// Server представляет HTTP сервер
type Server struct {
	config      *config.Config
	fileHandler *handlers.FileHandler
}

// New создает новый сервер
func New(cfg *config.Config, fileHandler *handlers.FileHandler) *Server {
	return &Server{
		config:      cfg,
		fileHandler: fileHandler,
	}
}

// Router настраивает и возвращает HTTP роутер
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	// Маршруты для работы с файлами
	mux.HandleFunc("/upload", s.fileHandler.UploadFile)
	mux.HandleFunc("/files", s.fileHandler.ListFiles)
	mux.HandleFunc("/download/", s.fileHandler.DownloadFile)

	// Статические файлы
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	// Главная страница
	mux.HandleFunc("/", s.handleHome)

	// Применяем middleware
	handler := middleware.CORS(mux)
	handler = middleware.Logging(handler)

	return handler
}

// handleHome обрабатывает главную страницу
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Файловый сервер</title>
    <meta charset="utf-8">
</head>
<body>
    <h1>Файловый сервер</h1>
    
    <h2>Загрузить файл</h2>
    <form action="/upload" method="post" enctype="multipart/form-data">
        <input type="file" name="file" required>
        <button type="submit">Загрузить</button>
    </form>
    
    <h2>Список файлов</h2>
    <div id="files"></div>
    
    <script>
        // Загружаем список файлов
        fetch('/files')
            .then(response => response.json())
            .then(files => {
                const filesDiv = document.getElementById('files');
                if (files.length === 0) {
                    filesDiv.innerHTML = '<p>Нет загруженных файлов</p>';
                    return;
                }
                
                const list = files.map(file => 
                    '<li><a href="/download/' + file.name + '">' + file.name + '</a> (' + file.size + ' байт)</li>'
                ).join('');
                
                filesDiv.innerHTML = '<ul>' + list + '</ul>';
            });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}