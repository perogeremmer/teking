package main

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ayomendaki/ayomendaki-admin/internal/database"
	"github.com/ayomendaki/ayomendaki-admin/internal/server"
)

//go:embed web/templates web/static
var assets embed.FS

func main() {
	dev := flag.Bool("dev", false, "development mode: read templates from disk")
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/ayomendaki.db"
	}

	if err := database.Init(dbPath); err != nil {
		log.Fatal("Database init failed:", err)
	}
	defer database.Close()

	if err := database.Seed(); err != nil {
		log.Fatal("Seed failed:", err)
	}

	mode := "production"
	if *dev {
		mode = "development"
	}
	log.Printf("Server starting on :%s (%s)", port, mode)

	srv := server.New(assets, *dev)
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatal("Server failed:", err)
	}
}
