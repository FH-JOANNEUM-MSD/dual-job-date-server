package routes

import (
	"net/http"

	"dual-job-date-server/internal/auth"
	"dual-job-date-server/internal/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// ==========================================
	// 🟢 ÖFFENTLICHE ROUTEN (Ohne Token)
	// ==========================================
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server läuft!"))
	}).Methods("GET")

	r.HandleFunc("/api/resend-invite", handlers.HandleResendInvite).Methods("POST")

	// ==========================================
	// 🔴 GESCHÜTZTE ROUTEN (Mit Token & Rollen)
	// ==========================================
	api := r.PathPrefix("/api").Subrouter()

	// JWT Middleware für alle geschützten Routen
	api.Use(auth.JWTMiddleware)

	// --- HELPER / SYSTEM ---
	api.HandleFunc("/seed", auth.RequireRole("admin")(handlers.SeedDatabase)).Methods("GET")
	api.HandleFunc("/invite", auth.RequireRole("admin")(handlers.InviteUserHandler)).Methods("POST")
	api.HandleFunc("/meetings/assign", auth.RequireRole("admin")(handlers.AssignMeetingsByPreferencesHandler)).Methods("POST")

	// --- COMPANIES ---
	// Listen & Interaktion
	api.HandleFunc("/companies/active", auth.RequireRole("admin", "student")(handlers.GetActiveCompaniesHandler)).Methods("GET")
	api.HandleFunc("/companies/{id}/vote", auth.RequireRole("student")(handlers.VoteCompanyHandler)).Methods("POST")

	// Self-Service (Entity-Typ: "company")
	api.HandleFunc("/companies/{id}", auth.RequireSelfOrAdmin("company")(handlers.UpdateCompanyHandler)).Methods("PATCH")
	api.HandleFunc("/companies/{id}/logo", auth.RequireSelfOrAdmin("company")(handlers.UploadCompanyLogoHandler)).Methods("POST")
	api.HandleFunc("/companies/{id}/images", auth.RequireSelfOrAdmin("company")(handlers.UploadCompanyImageHandler)).Methods("POST")

	// --- STUDENTS ---
	// Listen
	api.HandleFunc("/students", auth.RequireRole("admin")(handlers.GetAllStudentsHandler)).Methods("GET")

	// Self-Service (Entity-Typ: "student")
	api.HandleFunc("/students/{id}/preferences", auth.RequireSelfOrAdmin("student")(handlers.GetPreferencesByStudentHandler)).Methods("GET")
	api.HandleFunc("/students/{id}/meetings", auth.RequireSelfOrAdmin("student")(handlers.GetMeetingsByStudentHandler)).Methods("GET")
	api.HandleFunc("/students/{id}", auth.RequireSelfOrAdmin("student")(handlers.UpdateStudentHandler)).Methods("PATCH")
	api.HandleFunc("/students/{id}", auth.RequireSelfOrAdmin("student")(handlers.DeleteStudentHandler)).Methods("DELETE")

	// --- USERS (Generell) ---
	// Self-Service (Entity-Typ: "user")
	api.HandleFunc("/users/{id}", auth.RequireSelfOrAdmin("user")(handlers.UpdateUserNamesHandler)).Methods("PATCH")

	// --- GENERAL / INFO (Für alle eingeloggten User) ---
	api.HandleFunc("/events/active", handlers.GetActiveEventHandler).Methods("GET")
	api.HandleFunc("/slots", handlers.GetAllSlotsHandler).Methods("GET")
	api.HandleFunc("/me", handlers.GetMyIDHandler).Methods("GET")

	// Nur Admin darf Slots löschen
	api.HandleFunc("/slots/{id}", auth.RequireRole("admin")(handlers.DeleteSlotHandler)).Methods("DELETE")

	return r
}
