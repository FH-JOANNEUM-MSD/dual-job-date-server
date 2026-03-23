package handlers

import (
	"dual-job-date-server/internal/models"
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux" // <-- NEU: Gorilla Mux importiert
)

func UpdateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. ID sicher über Gorilla Mux aus der URL ({id}) holen
	vars := mux.Vars(r)
	idStr := vars["id"]

	companyID, err := strconv.Atoi(idStr)
	if err != nil {
		// Etwas detailliertere Fehlermeldung, falls mal wieder was schiefgeht
		http.Error(w, "Ungültige Firmen-ID: '"+idStr+"'", http.StatusBadRequest)
		return
	}

	// 2. Das JSON aus dem Request-Body in unser Pointer-Struct übersetzen
	var input models.UpdateCompanyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Fehlerhaftes JSON-Format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Das Repository aufrufen und die Daten übergeben
	err = repository.UpdateCompany(companyID, input)
	if err != nil {
		http.Error(w, "Fehler beim Updaten der Firma: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Erfolgsmeldung zurückgeben
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Firma erfolgreich aktualisiert",
		"status":  "success",
	})
}
