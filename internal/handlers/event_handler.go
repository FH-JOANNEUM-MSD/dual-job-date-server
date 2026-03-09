package handlers

import (
    "encoding/json"
    "net/http"
    "dual-job-date-server/internal/repository"
)

func GetActiveEventHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    event, err := repository.GetActiveEvent()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(event)
}