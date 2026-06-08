package handlers

import (
	"encoding/json"
	"net/http"

	"dual-job-date-server/internal/models"
	"dual-job-date-server/internal/repository"
)

func BulkImportHandler(w http.ResponseWriter, r *http.Request) {
	var input models.BulkImportRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Ungültiges JSON", http.StatusBadRequest)
		return
	}

	if len(input.Students) == 0 && len(input.Companies) == 0 {
		http.Error(w, "mindestens ein student oder eine company muss angegeben werden", http.StatusBadRequest)
		return
	}

	result, err := repository.BulkImportStudentsAndCompanies(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
}
