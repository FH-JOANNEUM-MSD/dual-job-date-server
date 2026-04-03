package handlers

import (
	"dual-job-date-server/internal/models"
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"net/http"
)

func HandleResendInvite(w http.ResponseWriter, r *http.Request) {
	var req models.InviteRequest

	// JSON parsen
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// WICHTIG: Wir rufen ReinviteAuthUser auf.
	// Wir fangen ID (_) und Error (err) ab -> Assignment Mismatch gelöst!
	_, err := repository.ReinviteAuthUser(req.Email, req.Role)
	if err != nil {
		http.Error(w, "Fehler beim Resend: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Einladung wurde erfolgreich erneut versendet.",
	})
}
