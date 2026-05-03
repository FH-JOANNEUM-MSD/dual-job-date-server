package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"dual-job-date-server/internal/repository"

	"github.com/gorilla/mux"
)

func GetCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		http.Error(w, "Missing company ID", http.StatusBadRequest)
		return
	}

	// String sicher in einen Integer umwandeln
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	// 🟢 HIER IST DIE ÄNDERUNG: Wir rufen den neuen Namen auf!
	company, err := repository.GetSingleCompanyByID(id)

	if err != nil {
		if err == repository.ErrCompanyNotFound {
			http.Error(w, "Company not found", http.StatusNotFound)
			return
		}

		// 🟢 HIER: Wir loggen den ECHTEN Fehler in deine Konsole!
		fmt.Printf("--- 🔴 DB FEHLER BEI COMPANY %d: %v ---\n", id, err)

		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}
