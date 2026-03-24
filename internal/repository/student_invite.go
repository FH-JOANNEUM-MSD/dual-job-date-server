package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"

	"github.com/google/uuid"
)

func CreateStudentProfile(authUUID string, req models.InviteRequest) error {
	// 1. Unsere interne ID generieren
	internalUUID := uuid.New().String()

	// 2. Den User in public.users eintragen
	userInsert := map[string]interface{}{
		"id":      internalUUID,
		"user_id": authUUID, // Die ID, die wir von InviteAuthUser bekommen haben
		"role":    "student",
	}
	err := database.SupabaseClient.DB.From("users").Insert(userInsert).Execute(nil)
	if err != nil {
		return err
	}

	// 3. Den Studenten in public.students eintragen
	studentInsert := map[string]interface{}{
		"user_id":       internalUUID, // Verknüpfung zu unserer internen users-Tabelle
		"study_program": req.StudyProgram,
		"semester":      req.Semester,
	}
	err = database.SupabaseClient.DB.From("students").Insert(studentInsert).Execute(nil)

	return err
}
