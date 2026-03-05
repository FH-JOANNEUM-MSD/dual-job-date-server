package main

import (

	"dual-job-date-server/internal/routes"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Lade .env Datei nur, wenn sie existiert
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Warnung: .env Datei konnte nicht geladen werden")
		}
	}

	r := routes.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:" + os.Getenv("API_PORT"),
		WriteTimeout: 80 * time.Second,
		ReadTimeout:  80 * time.Second,
		IdleTimeout:  50 * time.Second,
	}


	log.Println("Server läuft auf https://0.0.0.0:" + os.Getenv("API_PORT"))
	log.Fatal(srv.ListenAndServe())
}
