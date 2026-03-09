package repository

import (
    "dual-job-date-server/internal/database"
    "dual-job-date-server/internal/models"
    "strconv"
)

func GetPreferencesByStudent(studentID int) ([]models.Preference, error) {
    var preferences []models.Preference

    err := database.SupabaseClient.DB.From("preferences").Select("*").Eq("student_id", strconv.Itoa(studentID)).Execute(&preferences)
    if err != nil {
        return nil, err
    }

    return preferences, nil
}