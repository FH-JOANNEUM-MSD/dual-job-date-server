package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

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

func CreateMeetingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input models.CreateMeetingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Fehlerhaftes JSON-Format: "+err.Error(), http.StatusBadRequest)
		return
	}

	created, err := repository.CreateMeeting(input)
	if err != nil {
		if errors.Is(err, repository.ErrMeetingConflict) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if isMeetingReferenceError(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Fehler beim Anlegen des Meetings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func DeleteMeetingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	meetingID, err := strconv.Atoi(vars["id"])
	if err != nil || meetingID <= 0 {
		http.Error(w, "Ungültige Meeting-ID", http.StatusBadRequest)
		return
	}

	if err := repository.DeleteMeeting(meetingID); err != nil {
		if errors.Is(err, repository.ErrMeetingNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Fehler beim Löschen des Meetings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Meeting erfolgreich gelöscht",
		"status":  "success",
	})
}

func UpdateMeetingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	meetingID, err := strconv.Atoi(vars["id"])
	if err != nil || meetingID <= 0 {
		http.Error(w, "Ungültige Meeting-ID", http.StatusBadRequest)
		return
	}

	var input models.UpdateMeetingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Fehlerhaftes JSON-Format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if input.SlotID == nil && input.StudentID == nil && input.CompanyID == nil {
		http.Error(w, "Mindestens eines der Felder slot_id, student_id oder company_id muss gesetzt sein", http.StatusBadRequest)
		return
	}

	updated, err := repository.UpdateMeeting(meetingID, input)
	if err != nil {
		if errors.Is(err, repository.ErrMeetingNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrMeetingConflict) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if isMeetingReferenceError(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Fehler beim Updaten des Meetings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

func isMeetingReferenceError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "nicht gefunden") || strings.Contains(msg, "muss größer als 0 sein")
}
