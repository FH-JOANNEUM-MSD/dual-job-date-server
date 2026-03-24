package handlers

import (
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ungültige Studenten-ID: '"+idStr+"'", http.StatusBadRequest)
		return
	}

	// Die neue All-in-One Löschfunktion aufrufen
	err = repository.HardDeleteStudentAndUser(studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Wenn wir hier ankommen, ist der User komplett pulverisiert.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Student, alle Daten und Auth-Account restlos gelöscht",
		"status":  "success",
	})
}
