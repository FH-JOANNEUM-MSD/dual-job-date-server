package handlers

import (
    "encoding/json"
    "net/http"

    "dual-job-date-server/internal/models"
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

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var input models.CreateEventInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Fehlerhaftes JSON-Format: "+err.Error(), http.StatusBadRequest)
        return
    }

    if input.Name == "" || input.EventDate == "" {
        http.Error(w, "Felder 'name' und 'event_date' sind erforderlich", http.StatusBadRequest)
        return
    }

    event, err := repository.CreateEvent(input)
    if err != nil {
        http.Error(w, "Fehler beim Anlegen des Events: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(event)
}