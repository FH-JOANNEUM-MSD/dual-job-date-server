package auth

import (
	"context" // Wichtig: Neu dabei für den Datentransport
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Authorization Header auslesen
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Login erforderlich (Header fehlt)", http.StatusUnauthorized)
			return
		}

		// 2. "Bearer " Präfix entfernen
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Token validieren
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("Unerwarteter Algorithmus: %v", token.Header["alg"])
			}
			return GetSupabasePublicKey(), nil
		})

		// 4. Fehlerbehandlung
		if err != nil || !token.Valid {
			http.Error(w, "Ungültiger oder abgelaufener Token", http.StatusUnauthorized)
			return
		}

		// --- NEU: User-ID extrahieren und in den Request-"Rucksack" packen ---
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Supabase speichert die User-ID im Feld "sub"
			userID, ok := claims["sub"].(string)
			if ok {
				// Wir speichern die ID im Context des Requests
				ctx := context.WithValue(r.Context(), "userID", userID)
				r = r.WithContext(ctx)
			}
		}

		// 5. Alles okay? Dann weiter mit dem aktualisierten Request (inkl. UserID)
		next.ServeHTTP(w, r)
	})
}
