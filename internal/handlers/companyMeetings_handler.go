package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"dual-job-date-server/internal/repository"

	"github.com/gorilla/mux" // Wichtig: mux importieren!
)

func GetMeetingsByCompanyIDHandler(w http.ResponseWriter, r *http.Request) {
	// Wenn du .Methods("GET") im Router hast, ist dieser Check eigentlich
	// optional, aber schadet nicht zur doppelten Absicherung.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Die ID direkt aus den Pfad-Variablen von mux auslesen
	vars := mux.Vars(r)
	companyIDStr := vars["id"]

	if companyIDStr == "" {
		http.Error(w, "Company ID fehlt in der URL", http.StatusBadRequest)
		return
	}

	// 2. String zu Integer parsen
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		http.Error(w, "Ungültige Company ID", http.StatusBadRequest)
		return
	}

	// 3. Repository aufrufen
	meetings, err := repository.GetMeetingsByCompanyID(companyID)
	if err != nil {
		http.Error(w, "Fehler beim Laden der Meetings", http.StatusInternalServerError)
		return
	}

	// 4. JSON Response senden
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(meetings); err != nil {
		http.Error(w, "Fehler beim Parsen zu JSON", http.StatusInternalServerError)
	}
}
