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