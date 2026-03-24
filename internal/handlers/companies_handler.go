package handlers

import (
	"dual-job-date-server/internal/models"
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"net/http"
)

func GetActiveCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	companies, err := repository.GetActiveCompanies()
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
