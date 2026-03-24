package handlers

import (
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteSlotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. ID aus der URL holen (/api/slots/3)
	vars := mux.Vars(r)
	idStr := vars["id"]

	slotID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ungültige Slot-ID: '"+idStr+"'", http.StatusBadRequest)
		return
	}

	// 2. Das Repository aufrufen (erst Meetings löschen, dann den Slot)
	err = repository.DeleteSlot(slotID)
	if err != nil {
		http.Error(w, "Fehler beim Löschen des Slots: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Erfolgsmeldung zurückgeben
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Slot und zugehörige Meetings erfolgreich gelöscht",
		"status":  "success",
	})
}
