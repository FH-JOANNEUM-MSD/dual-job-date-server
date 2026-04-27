package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv" // 🟢 NEU: Wichtig für das sichere Parsen der ID
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

		var userResult []map[string]interface{}
		dbErr := database.SupabaseClient.DB.From("users").Select("id,role").Eq("user_id", userID).Execute(&userResult)

		if dbErr != nil {
			http.Error(w, "Datenbankfehler bei Rollenprüfung", http.StatusInternalServerError)
			return
		}

		if len(userResult) == 0 {
			http.Error(w, "User in der Datenbank nicht gefunden", http.StatusUnauthorized)
			return
		}

		role, roleOk := userResult[0]["role"].(string)
		dbUserID, idOk := userResult[0]["id"].(string)

		if !roleOk || !idOk {
			http.Error(w, "User-Datenbankeintrag unvollständig", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "auth_user_id", userID)
		ctx = context.WithValue(ctx, "db_user_id", dbUserID)
		ctx = context.WithValue(ctx, "role", role)

		// ==========================================
		// 🟢 FIX: Sicheres Parsen der IDs
		// ==========================================
		if role == "student" {
			var studentResult []map[string]interface{}
			// TIPP: Falls 'students' auf die neue users-Tabelle verweist, ändere 'userID' in 'dbUserID'
			err := database.SupabaseClient.DB.From("students").Select("id").Eq("user_id", dbUserID).Execute(&studentResult)

			if err == nil && len(studentResult) > 0 {
				var sID int
				// Robustes Abfangen: Egal ob Supabase Float oder String liefert!
				switch v := studentResult[0]["id"].(type) {
				case float64:
					sID = int(v)
				case string:
					sID, _ = strconv.Atoi(v)
				}
				if sID != 0 {
					ctx = context.WithValue(ctx, "student_id", sID)
				}
			}
		} else if role == "company" {
			var companyResult []map[string]interface{}
			err := database.SupabaseClient.DB.From("companies").Select("id").Eq("user_id", dbUserID).Execute(&companyResult)

			if err == nil && len(companyResult) > 0 {
				var cID int
				switch v := companyResult[0]["id"].(type) {
				case float64:
					cID = int(v)
				case string:
					cID, _ = strconv.Atoi(v)
				}
				if cID != 0 {
					ctx = context.WithValue(ctx, "company_id", cID)
				}
			}
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
