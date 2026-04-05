package handlers

import (
	"dual-job-date-server/internal/models"
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func UpdateUserNamesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. ID sicher über Gorilla Mux aus der URL ({id}) holen
	vars := mux.Vars(r)
	idStr := vars["id"]

	// ACHTUNG: Users nutzen UUIDs als Primary Key, keine Integers!
	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Ungültige User-ID Format: '"+idStr+"'", http.StatusBadRequest)
		return
	}

	// 2. Das JSON aus dem Request-Body in unser Pointer-Struct übersetzen
	var input models.UpdateUserNameInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Fehlerhaftes JSON-Format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Optionaler Check: Wenn beide Werte im JSON fehlen, können wir uns den DB-Call sparen
	if input.FirstName == nil && input.LastName == nil {
		http.Error(w, "Mindestens first_name oder last_name muss im JSON angegeben werden", http.StatusBadRequest)
		return
	}

	// 3. Das Repository aufrufen und die Daten übergeben
	err = repository.UpdateUserNames(userID, input)
	if err != nil {
		http.Error(w, "Fehler beim Updaten des Benutzers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Erfolgsmeldung zurückgeben
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Benutzer erfolgreich aktualisiert",
		"status":  "success",
	})
}
