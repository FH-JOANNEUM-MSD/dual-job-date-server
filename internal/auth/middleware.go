package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"dual-job-date-server/internal/database" // Dein globaler DB Client!

	"github.com/golang-jwt/jwt/v5"
)

// Die Signatur bleibt ganz normal, kein DB-Parameter nötig!
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Login erforderlich", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("Unerwarteter Algorithmus: %v", token.Header["alg"])
			}
			return GetSupabasePublicKey(), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Ungültiger Token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Fehler beim Lesen der Claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "Keine User-ID im Token", http.StatusUnauthorized)
			return
		}

		// --- NEU: Supabase Abfrage direkt hier! ---
		var userResult []map[string]interface{}

		// Select("id, role") statt nur Select("role")
		dbErr := database.SupabaseClient.DB.From("users").Select("id,role").Eq("user_id", userID).Execute(&userResult)

		if dbErr != nil {
			http.Error(w, "Datenbankfehler bei Rollenprüfung", http.StatusInternalServerError)
			return
		}

		if len(userResult) == 0 {
			http.Error(w, "User in der Datenbank nicht gefunden", http.StatusUnauthorized)
			return
		}

		// Rolle und DB-ID aus dem Ergebnis extrahieren
		role, roleOk := userResult[0]["role"].(string)
		dbUserID, idOk := userResult[0]["id"].(string) // <-- HIER IST DIE ECHTE DB-UUID!

		if !roleOk || !idOk {
			http.Error(w, "User-Datenbankeintrag unvollständig", http.StatusInternalServerError)
			return
		}

		// 1. Hier nutzen wir := (weil ctx neu erschaffen wird) und r.Context() als Basis
		ctx := context.WithValue(r.Context(), "userID", userID)

		// 2. Ab hier nutzen wir NUR NOCH = (weil ctx schon existiert) und stapeln auf dem aktuellen ctx!
		ctx = context.WithValue(ctx, "auth_user_id", userID)
		ctx = context.WithValue(ctx, "db_user_id", dbUserID)
		ctx = context.WithValue(ctx, "role", role)
		// ==========================================
		// NEU: Spezifische Integer-IDs (Student / Company) laden
		// ==========================================
		if role == "student" {
			var studentResult []map[string]interface{}
			err := database.SupabaseClient.DB.From("students").Select("id").Eq("user_id", userID).Execute(&studentResult)

			if err == nil && len(studentResult) > 0 {
				// WICHTIG: Supabase/JSON-Unmarshaling macht aus Zahlen standardmäßig float64!
				if sID, ok := studentResult[0]["id"].(float64); ok {
					ctx = context.WithValue(ctx, "student_id", int(sID))
				}
			}
		} else if role == "company" {
			var companyResult []map[string]interface{}
			err := database.SupabaseClient.DB.From("companies").Select("id").Eq("user_id", userID).Execute(&companyResult)

			if err == nil && len(companyResult) > 0 {
				if cID, ok := companyResult[0]["id"].(float64); ok {
					ctx = context.WithValue(ctx, "company_id", int(cID))
				}
			}
		}

		// Den aktualisierten Context an den Request anhängen
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
