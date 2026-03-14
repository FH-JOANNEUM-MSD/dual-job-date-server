package routes

import (
	"net/http"

	"dual-job-date-server/internal/auth"
	"dual-job-date-server/internal/handlers"

	"github.com/gorilla/mux"
)

// NewRouter erstellt den Router und registriert alle Routen
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// --- System & Test Routen ---
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server läuft!"))
	}).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(auth.JWTMiddleware)

	//Helper
	api.HandleFunc("/seed", handlers.SeedDatabase).Methods("GET")

	// --- API Routen (Business Logik) ---

	// Firmen-Endpunkte
	// Studenten sehen hier die Liste zum Voten
	api.HandleFunc("/companies/active", handlers.GetActiveCompaniesHandler).Methods("GET")

	// Studenten-Endpunkte
	// Hier können Admins alle Studierenden verwalten
	api.HandleFunc("/students", handlers.GetAllStudentsHandler).Methods("GET")

	// Spezifische Daten für einen Studenten (Meetings & Preferences)
	api.HandleFunc("/students/{id}/preferences", handlers.GetPreferencesByStudentHandler).Methods("GET")
	api.HandleFunc("/students/{id}/meetings", handlers.GetMeetingsByStudentHandler).Methods("GET")

	// Event & Zeitplan
	api.HandleFunc("/events/active", handlers.GetActiveEventHandler).Methods("GET")
	api.HandleFunc("/slots", handlers.GetAllSlotsHandler).Methods("GET")

	return r
}
