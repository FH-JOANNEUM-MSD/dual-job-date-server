package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"dual-job-date-server/internal/models"
	"dual-job-date-server/internal/repository"

	"github.com/gorilla/mux"
)

func GetMeetingsByStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	studentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Ungültige Studenten-ID", http.StatusBadRequest)
		return
	}

	meetings, err := repository.GetMeetingsByStudent(studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(meetings)
}

func GetMeetingsByCompanyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	companyID, err := strconv.Atoi(vars["id"])
	if err != nil || companyID <= 0 {
		http.Error(w, "Ungueltige Company-ID", http.StatusBadRequest)
		return
	}

	meetings, err := repository.GetMeetingsByCompanyWithStudentAndSlot(companyID)
	if err != nil {
		http.Error(w, "Fehler beim Laden der Termine: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if meetings == nil {
		meetings = []models.CompanyMeeting{}
	}

	json.NewEncoder(w).Encode(meetings)
}
