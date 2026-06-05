package repository

import (
    "dual-job-date-server/internal/database"
    "dual-job-date-server/internal/models"
)

func GetAllSlots() ([]models.Slot, error) {
    var slots []models.Slot

    err := database.SupabaseClient.DB.From("slots").Select("*").Execute(&slots)
    if err != nil {
        return nil, err
    }

    return slots, nil
}

func CreateSlot(input models.CreateSlotInput) (models.Slot, error) {
    insertData := map[string]interface{}{
        "start_time": input.StartTime,
        "end_time":   input.EndTime,
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