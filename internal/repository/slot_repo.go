package repository

import (
    "errors"
    "strconv"

    "dual-job-date-server/internal/database"
    "dual-job-date-server/internal/models"
)

var ErrSlotNotFound = errors.New("slot nicht gefunden")

func GetAllSlots(eventID int) ([]models.Slot, error) {
    var slots []models.Slot

    query := database.SupabaseClient.DB.From("slots").Select("*")
    if eventID > 0 {
        if err := query.Eq("event_id", strconv.Itoa(eventID)).Execute(&slots); err != nil {
            return nil, err
        }
        return slots, nil
    }

    if err := query.Execute(&slots); err != nil {
        return nil, err
    }
    return slots, nil
}

func CreateSlot(input models.CreateSlotInput) (models.Slot, error) {
    insertData := map[string]interface{}{
        "start_time": input.StartTime,
        "end_time":   input.EndTime,
        "event_id":   input.EventID,
    }

    var created []models.Slot
    err := database.SupabaseClient.DB.From("slots").Insert(insertData).Execute(&created)
    if err != nil {
        return models.Slot{}, err
    }
    if len(created) == 0 {
        return models.Slot{}, nil
    }
    return created[0], nil
}

func GetSlotByID(slotID int) (models.Slot, error) {
    var slots []models.Slot
    err := database.SupabaseClient.DB.From("slots").Select("*").Eq("id", strconv.Itoa(slotID)).Execute(&slots)
    if err != nil {
        return models.Slot{}, err
    }
    if len(slots) == 0 {
        return models.Slot{}, ErrSlotNotFound
    }
    return slots[0], nil
}

func UpdateSlot(slotID int, input models.UpdateSlotInput) (models.Slot, error) {
    if _, err := GetSlotByID(slotID); err != nil {
        return models.Slot{}, err
    }

    updateData := make(map[string]interface{})
    if input.StartTime != nil {
        updateData["start_time"] = *input.StartTime
    }
    if input.EndTime != nil {
        updateData["end_time"] = *input.EndTime
    }

    var updated []models.Slot
    err := database.SupabaseClient.DB.From("slots").Update(updateData).Eq("id", strconv.Itoa(slotID)).Execute(&updated)
    if err != nil {
        return models.Slot{}, err
    }
    if len(updated) == 0 {
        return models.Slot{}, ErrSlotNotFound
    }
    return updated[0], nil
}