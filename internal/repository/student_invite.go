package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"

	"github.com/google/uuid"
)

func CreateStudentProfile(authUUID string, req models.InviteRequest) error {
	internalUUID := uuid.New().String()

	// 1. Eintrag in public.users
	userInsert := map[string]interface{}{
		"id":         internalUUID, // Das ist der PK
		"user_id":    authUUID,     // ID vom Supabase Auth
		"role":       "student",
		"first_name": req.FirstName,
		"last_name":  req.LastName,
	}

	err := database.SupabaseClient.DB.From("users").Insert(userInsert).Execute(nil)
	if err != nil {
		return err // Falls hier was schiefläuft, direkt raus
	}

	// 2. Eintrag in public.students
	studentInsert := map[string]interface{}{
		"user_id":       internalUUID, // Muss exakt die ID von oben sein
		"study_program": req.StudyProgram,
	}

	// Nur mitschicken, wenn vorhanden (bei 0 wird es in DB NULL)
	if req.Semester > 0 {
		studentInsert["semester"] = req.Semester
	}

	err = database.SupabaseClient.DB.From("students").Insert(studentInsert).Execute(nil)
	return err
}
