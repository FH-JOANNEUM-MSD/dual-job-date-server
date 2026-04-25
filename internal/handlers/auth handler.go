package handlers

import (
	"encoding/json"
	"net/http"

	"dual-job-date-server/internal/auth"
	"dual-job-date-server/internal/models"
	"dual-job-date-server/internal/repository" // Dein Repo importieren
)

func GetMyIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Auth-UUID aus der Middleware
	userID := auth.GetUserID(r.Context())

	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Keine valide Session gefunden"})
		return
	}

	// 1. Repo aufrufen (🟢 NEU: Wir fangen dbID auf)
	dbID, role, studentID, companyID, err := repository.GetUserAuthDetails(userID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 2. Response bauen (🟢 NEU: ID zuweisen)
	response := models.UserAuthResponse{
		ID:        dbID,   // Die echte Datenbank-UUID
		UserID:    userID, // Die Supabase-Auth-ID
		Status:    "authenticated",
		Role:      role,
		StudentID: studentID,
		CompanyID: companyID,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
