package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
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