package main

import (
	"log"
	"net/http"

	"file-server-go/internal/config"
	"file-server-go/internal/server"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Создаем и запускаем сервер
	srv := server.New(cfg)

	log.Printf("Сервер запущен на порту %s", cfg.Port)
	log.Printf("Папка для загрузок: %s", cfg.UploadDir)

	if err := http.ListenAndServe(":"+cfg.Port, srv.Router()); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
