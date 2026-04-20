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

		// Wir fragen direkt deine 'users' Tabelle über den globalen SupabaseClient ab
		dbErr := database.SupabaseClient.DB.From("users").Select("role").Eq("user_id", userID).Execute(&userResult)

		if dbErr != nil {
			http.Error(w, "Datenbankfehler bei Rollenprüfung", http.StatusInternalServerError)
			return
		}

		if len(userResult) == 0 {
			http.Error(w, "User in der Datenbank nicht gefunden", http.StatusUnauthorized)
			return
		}

		// Rolle aus dem Ergebnis extrahieren (als String casten)
		role, ok := userResult[0]["role"].(string)
		if !ok {
			http.Error(w, "Rolle konnte nicht gelesen werden", http.StatusInternalServerError)
			return
		}

		// 4. userID UND role in den Context packen
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "role", role)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
