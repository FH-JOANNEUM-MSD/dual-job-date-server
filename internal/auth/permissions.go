package auth

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ==========================================
// REQUIRE ROLE
// ==========================================
func RequireRole(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			role, ok := r.Context().Value("role").(string)
			if !ok || role == "" {
				http.Error(w, "Unauthorized: missing role", http.StatusUnauthorized)
				return
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Printf("[RBAC ROLE BLOCK] role=%s required=%v", role, allowedRoles)
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}
}

// ==========================================
// REQUIRE SELF OR ADMIN (CLEAN + SAFE)
// ==========================================
func RequireSelfOrAdmin(entityType string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			role, _ := r.Context().Value("role").(string)

			// 🟢 NEU: Wir laden die echte Datenbank-UUID aus dem Context!
			dbUserID, _ := r.Context().Value("db_user_id").(string)
			authUserID, _ := r.Context().Value("auth_user_id").(string) // Fürs saubere Logging

			studentID, _ := r.Context().Value("student_id").(int)
			companyID, _ := r.Context().Value("company_id").(int)

			// ==========================================
			// ADMIN = ALWAYS PASS
			// ==========================================
			if role == "admin" {
				next.ServeHTTP(w, r)
				return
			}

			vars := mux.Vars(r)
			targetID := vars["id"]

			// ==========================================
			// USER (UUID)
			// ==========================================
			if entityType == "user" {
				// 🟢 NEU: Vergleicht die URL-ID mit der Datenbank-ID
				if targetID == dbUserID {
					next.ServeHTTP(w, r)
					return
				}

				log.Printf("[BLOCK USER] dbUser=%s target=%s (AuthID war: %s)", dbUserID, targetID, authUserID)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// ==========================================
			// STUDENT (INT SAFE PARSE)
			// ==========================================
			if entityType == "student" {

				if role != "student" {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}

				targetInt, err := strconv.Atoi(targetID)
				if err != nil {
					log.Printf("[BAD REQUEST] invalid student id: %s", targetID)
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}

				if targetInt == studentID {
					next.ServeHTTP(w, r)
					return
				}

				log.Printf("[BLOCK STUDENT] student=%d target=%d", studentID, targetInt)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// ==========================================
			// COMPANY (INT SAFE PARSE)
			// ==========================================
			if entityType == "company" {

				if role != "company" {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}

				targetInt, err := strconv.Atoi(targetID)
				if err != nil {
					log.Printf("[BAD REQUEST] invalid company id: %s", targetID)
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}

				if targetInt == companyID {
					next.ServeHTTP(w, r)
					return
				}

				log.Printf("[BLOCK COMPANY] company=%d target=%d", companyID, targetInt)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// ==========================================
			// DEFAULT BLOCK
			// ==========================================
			log.Printf("[RBAC BLOCK] role=%s dbUserID=%s type=%s target=%s",
				role, dbUserID, entityType, targetID)

			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}
}
