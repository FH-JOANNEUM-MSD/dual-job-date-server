package repository

import (
    "dual-job-date-server/internal/database"
    "dual-job-date-server/internal/models"
)

func GetActiveEvent() (models.Event, error) {
    var events []models.Event

    // Wir lassen .Limit(1) weg, da die Library es an dieser Stelle nicht unterstützt
    err := database.SupabaseClient.DB.From("events").Select("*").Eq("is_active", "true").Execute(&events)
    if err != nil {
        return models.Event{}, err
    }

    // Wenn die Liste leer ist, geben wir ein leeres Modell zurück
    if len(events) == 0 {
        return models.Event{}, nil 
    }

    // Wir geben einfach das erste gefundene aktive Event zurück
    return events[0], nil
}