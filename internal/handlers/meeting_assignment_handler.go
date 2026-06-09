package handlers

import (
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"io"
	"net/http"
)

type assignMeetingsRequest struct {
	EventID                  int   `json:"event_id"`
	SlotIDs                  []int `json:"slot_ids"`
	StudentIDs               []int `json:"student_ids"`
	DryRun                   bool  `json:"dry_run"`
	IncludeInactiveCompanies bool  `json:"include_inactive_companies"`
	ReplaceExisting          bool  `json:"replace_existing"`
}

func AssignMeetingsByPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req assignMeetingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
		http.Error(w, "Ungültiger JSON-Body", http.StatusBadRequest)
		return
	}

	if req.EventID <= 0 {
		http.Error(w, "Feld 'event_id' ist erforderlich", http.StatusBadRequest)
		return
	}

	result, err := repository.AssignMeetingsByPreferences(repository.AssignMeetingsOptions{
		EventID:                  req.EventID,
		SlotIDs:                  req.SlotIDs,
		StudentIDs:               req.StudentIDs,
		DryRun:                   req.DryRun,
		IncludeInactiveCompanies: req.IncludeInactiveCompanies,
		ReplaceExisting:          req.ReplaceExisting,
	})
	if err != nil {
		http.Error(w, "Fehler bei der Zuteilung: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !req.DryRun && result.InsertedMeetings > 0 {
		w.WriteHeader(http.StatusCreated)
	}

	json.NewEncoder(w).Encode(result)
}
