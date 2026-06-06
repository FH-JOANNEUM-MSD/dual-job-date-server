package repository

import (
    "errors"
    "strconv"

    "dual-job-date-server/internal/database"
    "dual-job-date-server/internal/models"
)

var ErrEventNotFound = errors.New("event nicht gefunden")

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

func GetEventByID(eventID int) (models.Event, error) {
    var events []models.Event
    err := database.SupabaseClient.DB.From("events").Select("*").Eq("id", strconv.Itoa(eventID)).Execute(&events)
    if err != nil {
        return models.Event{}, err
    }
    if len(events) == 0 {
        return models.Event{}, ErrEventNotFound
    }
    return events[0], nil
}

func UpdateEvent(eventID int, input models.UpdateEventInput) (models.Event, error) {
    if _, err := GetEventByID(eventID); err != nil {
        return models.Event{}, err
    }

    updateData := make(map[string]interface{})
    if input.Name != nil {
        updateData["name"] = *input.Name
    }
    if input.Location != nil {
        updateData["location"] = *input.Location
    }
    if input.Description != nil {
        updateData["description"] = *input.Description
    }
    if input.EventDate != nil {
        updateData["event_date"] = *input.EventDate
    }
    if input.IsActive != nil {
        updateData["is_active"] = *input.IsActive
    }

    var updated []models.Event
    err := database.SupabaseClient.DB.From("events").Update(updateData).Eq("id", strconv.Itoa(eventID)).Execute(&updated)
    if err != nil {
        return models.Event{}, err
    }
    if len(updated) == 0 {
        return models.Event{}, ErrEventNotFound
    }
    return updated[0], nil
}

func DeleteEvent(eventID int) error {
    var deleted []interface{}
    return database.SupabaseClient.DB.From("events").Delete().Eq("id", strconv.Itoa(eventID)).Execute(&deleted)
}

func CreateEvent(input models.CreateEventInput) (models.Event, error) {
    isActive := input.IsActive != nil && *input.IsActive

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