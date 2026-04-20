package routes

import (
	"net/http"

	"dual-job-date-server/internal/auth"
	"dual-job-date-server/internal/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server läuft!"))
	}).Methods("GET")

	r.HandleFunc("/api/resend-invite", handlers.HandleResendInvite).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	// JWT Middleware für alle geschützten Routen
	api.Use(auth.JWTMiddleware)

	// --- HELPER ---
	// Nur Admin darf die Datenbank seeden
	api.HandleFunc("/seed", auth.RequireRole("admin")(handlers.SeedDatabase)).Methods("GET")
	api.HandleFunc("/invite", auth.RequireRole("admin")(handlers.InviteUserHandler)).Methods("POST")

	// --- COMPANIES ---
	// Studenten & Admins dürfen die aktiven Firmen sehen
	api.HandleFunc("/companies/active", auth.RequireRole("admin", "student")(handlers.GetActiveCompaniesHandler)).Methods("GET")

	// Nur Studenten dürfen voten
	api.HandleFunc("/companies/{id}/vote", auth.RequireRole("student")(handlers.VoteCompanyHandler)).Methods("POST")

	// Admin oder Company selbst darf updaten/Logo hochladen
	api.HandleFunc("/companies/{id}", auth.RequireSelfOrAdmin()(handlers.UpdateCompanyHandler)).Methods("PATCH")
	api.HandleFunc("/companies/{id}/logo", auth.RequireSelfOrAdmin()(handlers.UploadCompanyLogoHandler)).Methods("POST")

	// --- STUDENTS ---
	// Nur Admin darf alle Studenten sehen
	api.HandleFunc("/students", auth.RequireRole("admin")(handlers.GetAllStudentsHandler)).Methods("GET")

	// Eigene Daten: Admin oder Student selbst
	api.HandleFunc("/students/{id}/preferences", auth.RequireSelfOrAdmin()(handlers.GetPreferencesByStudentHandler)).Methods("GET")
	api.HandleFunc("/students/{id}/meetings", auth.RequireSelfOrAdmin()(handlers.GetMeetingsByStudentHandler)).Methods("GET")
	api.HandleFunc("/students/{id}", auth.RequireSelfOrAdmin()(handlers.UpdateStudentHandler)).Methods("PATCH")
	api.HandleFunc("/students/{id}", auth.RequireSelfOrAdmin()(handlers.DeleteStudentHandler)).Methods("DELETE")

	// --- GENERAL / SYSTEM ---
	// Admin triggert den Matching-Algo
	api.HandleFunc("/meetings/assign", auth.RequireRole("admin")(handlers.AssignMeetingsByPreferencesHandler)).Methods("POST")

	// Alle eingeloggten User dürfen Events, Slots und sich selbst sehen
	api.HandleFunc("/events/active", handlers.GetActiveEventHandler).Methods("GET")
	api.HandleFunc("/slots", handlers.GetAllSlotsHandler).Methods("GET")
	api.HandleFunc("/me", handlers.GetMyIDHandler).Methods("GET")

	// Admin ODER User selbst dürfen Account updaten
	api.HandleFunc("/users/{id}", auth.RequireSelfOrAdmin()(handlers.UpdateUserNamesHandler)).Methods("PATCH")

	// Nur Admin darf Slots löschen
	api.HandleFunc("/slots/{id}", auth.RequireRole("admin")(handlers.DeleteSlotHandler)).Methods("DELETE")

	return r
}
