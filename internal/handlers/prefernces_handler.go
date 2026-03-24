package handlers

import (
	"dual-job-date-server/internal/auth"
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPreferencesByStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Holt die {id} aus der URL
	vars := mux.Vars(r)
	studentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Ungültige Studenten-ID", http.StatusBadRequest)
		return
	}

	preferences, err := repository.GetPreferencesByStudent(studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(preferences)
}

type VoteCompanyRequest struct {
	Vote string `json:"vote"`
}

func VoteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := auth.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, "Keine valide Session gefunden", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	companyID, err := strconv.Atoi(vars["id"])
	if err != nil || companyID <= 0 {
		http.Error(w, "Ungueltige Company-ID", http.StatusBadRequest)
		return
	}

	var req VoteCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungueltiger Request-Body", http.StatusBadRequest)
		return
	}

	savedPreference, err := repository.SaveVoteForStudentUser(userID, companyID, req.Vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(savedPreference)
}
