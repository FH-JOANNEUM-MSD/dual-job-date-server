package handlers

import (
	"encoding/json"
	"net/http"

	"dual-job-date-server/internal/repository"
)

func GetAllPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prefs, err := repository.GetAllPreferences()
	if err != nil {
		http.Error(w, "Fehler beim Laden der Preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(prefs); err != nil {
		http.Error(w, "Fehler beim Parsen zu JSON", http.StatusInternalServerError)
	}
}
