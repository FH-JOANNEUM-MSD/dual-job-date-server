package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"dual-job-date-server/internal/repository"

	"github.com/gorilla/mux"
)

type setEventMeetingsRequest struct {
	Meetings []struct {
		SlotID    int `json:"slot_id"`
		StudentID int `json:"student_id"`
		CompanyID int `json:"company_id"`
	} `json:"meetings"`
}

// SetEventMeetingsHandler bedient PUT /api/events/{id}/meetings (Admin): ersetzt den Zeitplan eines Events.
func SetEventMeetingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil || eventID <= 0 {
		http.Error(w, "Ungültige Event-ID", http.StatusBadRequest)
		return
	}

	var req setEventMeetingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Fehlerhaftes JSON-Format: "+err.Error(), http.StatusBadRequest)
		return
	}

	desired := make([]repository.AssignedMeeting, 0, len(req.Meetings))
	for _, m := range req.Meetings {
		if m.SlotID <= 0 || m.StudentID <= 0 || m.CompanyID <= 0 {
			http.Error(w, "Jedes Meeting braucht positive slot_id, student_id und company_id", http.StatusBadRequest)
			return
		}
		desired = append(desired, repository.AssignedMeeting{
			SlotID:    m.SlotID,
			StudentID: m.StudentID,
			CompanyID: m.CompanyID,
		})
	}

	meetings, err := repository.SetEventMeetings(eventID, desired)
	if err != nil {
		http.Error(w, "Fehler beim Speichern des Zeitplans: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(meetings)
}
