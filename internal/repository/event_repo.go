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

func GetAllEvents() ([]models.Event, error) {
    var events []models.Event
    err := database.SupabaseClient.DB.From("events").Select("*").Execute(&events)
    if err != nil {
        return nil, err
    }
    return events, nil
}

func deactivateAllEvents() error {
    var updated []models.Event
    return database.SupabaseClient.DB.From("events").
        Update(map[string]interface{}{"is_active": false}).
        Eq("is_active", "true").
        Execute(&updated)
}

func CreateEvent(input models.CreateEventInput) (models.Event, error) {
    isActive := input.IsActive != nil && *input.IsActive
    if isActive {
        if err := deactivateAllEvents(); err != nil {
            return models.Event{}, err
        }
    }

    insertData := map[string]interface{}{
        "name":        input.Name,
        "location":    input.Location,
        "description": input.Description,
        "event_date":  input.EventDate,
        "is_active":   isActive,
    }

    var created []models.Event
    err := database.SupabaseClient.DB.From("events").Insert(insertData).Execute(&created)
    if err != nil {
        return models.Event{}, err
    }
    if len(created) == 0 {
        return models.Event{}, nil
    }
    return created[0], nil
}