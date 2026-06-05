package handlers

import (
    "encoding/json"
    "net/http"

    "dual-job-date-server/internal/models"
    "dual-job-date-server/internal/repository"
)

func GetAllSlotsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    slots, err := repository.GetAllSlots()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(slots)
}

func CreateSlotHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var input models.CreateSlotInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Fehlerhaftes JSON-Format: "+err.Error(), http.StatusBadRequest)
        return
    }

    if input.StartTime == "" || input.EndTime == "" {
        http.Error(w, "Felder 'start_time' und 'end_time' sind erforderlich", http.StatusBadRequest)
        return
    }

    slot, err := repository.CreateSlot(input)
    if err != nil {
        http.Error(w, "Fehler beim Anlegen des Slots: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(slot)
}