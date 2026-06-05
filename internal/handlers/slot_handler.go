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

func UpdateSlotHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    vars := mux.Vars(r)
    slotID, err := strconv.Atoi(vars["id"])
    if err != nil || slotID <= 0 {
        http.Error(w, "Ungültige Slot-ID", http.StatusBadRequest)
        return
    }

    var input models.UpdateSlotInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Fehlerhaftes JSON-Format: "+err.Error(), http.StatusBadRequest)
        return
    }

    if input.StartTime == nil && input.EndTime == nil {
        http.Error(w, "Mindestens ein Feld (start_time, end_time) muss gesetzt sein", http.StatusBadRequest)
        return
    }

    slot, err := repository.UpdateSlot(slotID, input)
    if err != nil {
        if errors.Is(err, repository.ErrSlotNotFound) {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        http.Error(w, "Fehler beim Updaten des Slots: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(slot)
}