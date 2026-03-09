package handlers

import (
    "encoding/json"
    "net/http"
    "dual-job-date-server/internal/repository"
)

func GetAllStudentsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    students, err := repository.GetAllStudents()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(students)
}