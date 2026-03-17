// internal/handlers/auth_handler.go
package handlers

import (
	"dual-job-date-server/internal/auth" // Pfad anpassen!
	"dual-job-date-server/internal/models"
	"encoding/json"
	"net/http"
)

// GetMyIDHandler gibt die ID des aktuell eingeloggten Users zurück
func GetMyIDHandler(w http.ResponseWriter, r *http.Request) {
	// Hol die ID aus dem Context (die deine Middleware dort abgelegt hat)
	userID := auth.GetUserID(r.Context())

	if userID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Keine valide Session gefunden"})
		return
	}

	response := models.UserAuthResponse{
		UserID: userID,
		Status: "authenticated",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
