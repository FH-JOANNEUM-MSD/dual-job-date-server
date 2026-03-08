package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"dual-job-date-server/internal/handlers"
)

// NewRouter erstellt den Router und registriert alle Routen
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server läuft!"))
	}).Methods("GET")

	r.HandleFunc("/seed", handlers.SeedDatabase).Methods("GET")

	return r
}