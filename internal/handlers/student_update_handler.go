package handlers

import (
	"dual-job-date-server/internal/models"
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Ungültige Studenten-ID: '"+idStr+"'", http.StatusBadRequest)
		return
	}

	var input models.UpdateStudentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Fehlerhaftes JSON-Format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if input.StudyProgram == nil && input.Semester == nil && input.FirstName == nil && input.LastName == nil {
		http.Error(w, "Mindestens eines der Felder study_program, semester, first_name oder last_name muss gesetzt sein", http.StatusBadRequest)
		return
	}

	err = repository.UpdateStudent(studentID, input)
	if err != nil {
		if err == repository.ErrStudentNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Fehler beim Updaten des Studenten: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Student erfolgreich aktualisiert",
		"status":  "success",
	})
}
