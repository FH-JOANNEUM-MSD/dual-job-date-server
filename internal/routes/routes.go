package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter erstellt den Router und registriert alle Routen
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server läuft!"))
	}).Methods("GET")



	return r
}