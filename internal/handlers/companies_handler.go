package handlers

import (
	"dual-job-date-server/internal/auth"
	"dual-job-date-server/internal/models"
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetActiveCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	onlyUnvoted := false
	onlyUnvotedRaw := r.URL.Query().Get("only_unvoted")
	if onlyUnvotedRaw != "" {
		parsed, err := strconv.ParseBool(onlyUnvotedRaw)
		if err != nil {
			http.Error(w, "Ungueltiger Query-Parameter only_unvoted (erlaubt: true/false)", http.StatusBadRequest)
			return
		}
		onlyUnvoted = parsed
	}

	userID := auth.GetUserID(r.Context())
	companies, err := repository.GetActiveCompaniesForStudentUser(userID, onlyUnvoted)
	if err != nil {
		http.Error(w, "Fehler beim Laden der Unternehmen: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Hier wurde "models" verwendet, ohne dass es oben importiert war
	if companies == nil {
		companies = []models.Company{}
	}

	json.NewEncoder(w).Encode(companies)
}
