package repository

import (
    "dual-job-date-server/internal/database"
    "dual-job-date-server/internal/models"
    "strconv"
)

func GetMeetingsByStudent(studentID int) ([]models.Meeting, error) {
    var meetings []models.Meeting

    err := database.SupabaseClient.DB.From("meetings").Select("*").Eq("student_id", strconv.Itoa(studentID)).Execute(&meetings)
    if err != nil {
        return nil, err
    }

    return meetings, nil
}