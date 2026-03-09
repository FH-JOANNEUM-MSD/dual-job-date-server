package routes

import (
    "net/http"

    "github.com/gorilla/mux"
    "dual-job-date-server/internal/handlers"
)

// NewRouter erstellt den Router und registriert alle Routen
func NewRouter() *mux.Router {
    r := mux.NewRouter()

    // --- System & Test Routen ---
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Server läuft!"))
    }).Methods("GET")

	//Helper
    r.HandleFunc("/seed", handlers.SeedDatabase).Methods("GET")

    // --- API Routen (Business Logik) ---
    
    // Firmen-Endpunkte
    // Studenten sehen hier die Liste zum Voten
    r.HandleFunc("/api/companies/active", handlers.GetActiveCompaniesHandler).Methods("GET")

    // Studenten-Endpunkte
    // Hier können Admins alle Studierenden verwalten
    r.HandleFunc("/api/students", handlers.GetAllStudentsHandler).Methods("GET")
    
    // Spezifische Daten für einen Studenten (Meetings & Preferences)
    r.HandleFunc("/api/students/{id}/preferences", handlers.GetPreferencesByStudentHandler).Methods("GET")
    r.HandleFunc("/api/students/{id}/meetings", handlers.GetMeetingsByStudentHandler).Methods("GET")
    
    // Event & Zeitplan 
    r.HandleFunc("/api/events/active", handlers.GetActiveEventHandler).Methods("GET")
    r.HandleFunc("/api/slots", handlers.GetAllSlotsHandler).Methods("GET")

    return r
}