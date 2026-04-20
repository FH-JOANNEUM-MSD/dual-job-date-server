package auth

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RequireRole prüft, ob der User eine der erlaubten Rollen hat (z.B. nur Admin)
func RequireRole(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Hole die Rolle aus dem JWT-Context (muss deine JWTMiddleware setzen!)
			userRole, ok := r.Context().Value("role").(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check ob die Rolle in der Liste der erlaubten Rollen ist
			for _, role := range allowedRoles {
				if userRole == role {
					next.ServeHTTP(w, r) // Erlaubt! Leite an den echten Handler weiter
					return
				}
			}

			// Keine passende Rolle gefunden
			http.Error(w, "Forbidden: Unzureichende Berechtigungen", http.StatusForbidden)
		}
	}
}

// RequireSelfOrAdmin prüft bei Routen mit {id}, ob der User Admin ist oder die ID seine eigene ist
func RequireSelfOrAdmin() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userRole, _ := r.Context().Value("role").(string)
			userID, _ := r.Context().Value("userID").(string) // Die ID aus dem Token

			vars := mux.Vars(r)
			targetID := vars["id"] // Die ID aus der URL

			// Admin darf alles, User darf seine eigenen Daten bearbeiten
			if userRole == "admin" || userID == targetID {
				next.ServeHTTP(w, r) // Erlaubt!
				return
			}

			http.Error(w, "Forbidden: Du darfst nur deine eigenen Daten bearbeiten", http.StatusForbidden)
		}
	}
}
