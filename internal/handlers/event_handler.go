package handlers

import (
    "encoding/json"
    "errors"
    "net/http"
    "strconv"

    "dual-job-date-server/internal/models"
    "dual-job-date-server/internal/repository"

    "github.com/gorilla/mux"
)

func GetActiveEventHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    events, err := repository.GetActiveEvents()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if events == nil {
        events = []models.Event{}
    }

    json.NewEncoder(w).Encode(events)
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

func GetAllEventsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    events, err := repository.GetAllEvents()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if events == nil {
        events = []models.Event{}
    }

    json.NewEncoder(w).Encode(events)
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    vars := mux.Vars(r)
    eventID, err := strconv.Atoi(vars["id"])
    if err != nil || eventID <= 0 {
        http.Error(w, "Ungültige Event-ID", http.StatusBadRequest)
        return
    }

    var input models.UpdateEventInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Fehlerhaftes JSON-Format: "+err.Error(), http.StatusBadRequest)
        return
    }

    if input.Name == nil && input.Location == nil && input.Description == nil && input.EventDate == nil && input.IsActive == nil {
        http.Error(w, "Mindestens ein Feld muss gesetzt sein", http.StatusBadRequest)
        return
    }

    event, err := repository.UpdateEvent(eventID, input)
    if err != nil {
        if errors.Is(err, repository.ErrEventNotFound) {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        http.Error(w, "Fehler beim Updaten des Events: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(event)
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    vars := mux.Vars(r)
    eventID, err := strconv.Atoi(vars["id"])
    if err != nil || eventID <= 0 {
        http.Error(w, "Ungültige Event-ID", http.StatusBadRequest)
        return
    }

    if err := repository.DeleteEvent(eventID); err != nil {
        http.Error(w, "Fehler beim Löschen des Events: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Event erfolgreich gelöscht",
        "status":  "success",
    })
}