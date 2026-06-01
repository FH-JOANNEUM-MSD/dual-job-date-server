package handlers

import (
	"encoding/json"
	"net/http"

	"dual-job-date-server/internal/models"
	"dual-job-date-server/internal/repository"
)

// GetAllMeetingsHandler bedient die /allMeetings Route für den Admin
func GetAllMeetingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	meetings, err := repository.GetAllMeetingsWithStudentAndSlot()
	if err != nil {
		http.Error(w, "Fehler beim Laden der Termine: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if meetings == nil {
		meetings = []models.CompanyMeeting{}
	}

	json.NewEncoder(w).Encode(meetings)
}
