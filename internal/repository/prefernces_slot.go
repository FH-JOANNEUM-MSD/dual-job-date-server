package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
	"fmt"
	"strconv"
	"strings"
)

func GetPreferencesByStudent(studentID int) ([]models.Preference, error) {
	var preferences []models.Preference

	err := database.SupabaseClient.DB.From("preferences").Select("*").Eq("student_id", strconv.Itoa(studentID)).Execute(&preferences)
	if err != nil {
		return nil, err
	}

	return preferences, nil
}

func SaveVoteForStudentUser(userID string, companyID int, vote string) (models.Preference, error) {
	if userID == "" {
		return models.Preference{}, fmt.Errorf("keine user id vorhanden")
	}
	if companyID <= 0 {
		return models.Preference{}, fmt.Errorf("ungueltige company id")
	}

	normalizedVote := strings.ToLower(strings.TrimSpace(vote))
	if normalizedVote != "like" && normalizedVote != "dislike" && normalizedVote != "neutral" {
		return models.Preference{}, fmt.Errorf("ungueltiger vote, erlaubt: like, dislike, neutral")
	}

	studentID, err := getStudentIDByUserID(userID)
	if err != nil {
		return models.Preference{}, err
	}

	company, err := GetCompanyByID(companyID)
	if err != nil {
		return models.Preference{}, err
	}
	if company.ID == 0 {
		return models.Preference{}, fmt.Errorf("company nicht gefunden")
	}
	if !company.Active {
		return models.Preference{}, fmt.Errorf("company ist nicht aktiv")
	}

	var existing []models.Preference
	err = database.SupabaseClient.DB.
		From("preferences").
		Select("*").
		Eq("student_id", strconv.Itoa(studentID)).
		Eq("company_id", strconv.Itoa(companyID)).
		Execute(&existing)
	if err != nil {
		return models.Preference{}, err
	}

	if len(existing) > 0 {
		update := map[string]interface{}{
			"preference_type": normalizedVote,
		}
		var updated []models.Preference
		err = database.SupabaseClient.DB.
			From("preferences").
			Update(update).
			Eq("id", strconv.Itoa(existing[0].ID)).
			Execute(&updated)
		if err != nil {
			return models.Preference{}, err
		}
		if len(updated) == 0 {
			existing[0].PreferenceType = normalizedVote
			return existing[0], nil
		}
		return updated[0], nil
	}

	insert := map[string]interface{}{
		"student_id":      studentID,
		"company_id":      companyID,
		"preference_type": normalizedVote,
	}

	var inserted []models.Preference
	err = database.SupabaseClient.DB.
		From("preferences").
		Insert(insert).
		Execute(&inserted)
	if err != nil {
		return models.Preference{}, err
	}
	if len(inserted) == 0 {
		return models.Preference{
			StudentID:      studentID,
			CompanyID:      companyID,
			PreferenceType: normalizedVote,
		}, nil
	}
	return inserted[0], nil
}

func getStudentIDByUserID(userID string) (int, error) {
	// JWT "sub" liefert die echte Supabase-Auth-UUID.
	// In diesem Projekt zeigt students.user_id aber auf users.id (interne UUID).
	// Deshalb: auth-uuid -> users.id -> students.id
	var users []map[string]interface{}
	err := database.SupabaseClient.DB.
		From("users").
		Select("id").
		Eq("user_id", userID).
		Execute(&users)
	if err != nil {
		return 0, err
	}
	if len(users) == 0 {
		return 0, fmt.Errorf("user nicht gefunden")
	}

	internalUserID, ok := users[0]["id"].(string)
	if !ok || internalUserID == "" {
		return 0, fmt.Errorf("ungueltige users.id")
	}

	var students []models.Student
	err = database.SupabaseClient.DB.
		From("students").
		Select("*").
		Eq("user_id", internalUserID).
		Execute(&students)
	if err != nil {
		return 0, err
	}
	if len(students) == 0 {
		return 0, fmt.Errorf("student fuer user nicht gefunden")
	}
	return students[0].ID, nil
}
