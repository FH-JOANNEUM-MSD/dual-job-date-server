package main

import (
	"dual-job-date-server/internal/routes"
	"dual-job-date-server/internal/database"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	// Lade .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warnung: .env konnte nicht geladen werden")
	}

	// DB init
	database.InitSupabase()

	r := routes.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:" + os.Getenv("API_PORT"),
		WriteTimeout: 80 * time.Second,
		ReadTimeout:  80 * time.Second,
		IdleTimeout:  50 * time.Second,
	}

	log.Println("Server läuft auf http://0.0.0.0:" + os.Getenv("API_PORT"))

	log.Fatal(srv.ListenAndServe())
}