package main

import (
	"log"
	"net/http"

	"file-server-go/internal/config"
	"file-server-go/internal/database"
	"file-server-go/internal/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create database connection
	db, err := database.New()
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	// Create and start server
	srv := server.New(cfg, nil, db)

	log.Printf("Server started on port %s", cfg.Port)
	log.Printf("Upload directory: %s", cfg.UploadDir)

	if err := http.ListenAndServe(":"+cfg.Port, srv.Router()); err != nil {
		log.Fatal("Server startup error:", err)
	}
}
