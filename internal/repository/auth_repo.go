package repository

import (
	"dual-job-date-server/internal/database"
	"fmt"
)

type userRow struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type profileRow struct {
	ID int `json:"id"`
}

func GetUserAuthDetails(authUUID string) (role string, studentID *int, companyID *int, err error) {
	var users []userRow

	// FIX: Kein Leerzeichen bei "id,role" !!!
	err = database.SupabaseClient.DB.From("users").Select("id,role").Eq("user_id", authUUID).Execute(&users)

	if err != nil {
		return "", nil, nil, fmt.Errorf("datenbankfehler bei user-abfrage: %v", err)
	}
	if len(users) == 0 {
		return "", nil, nil, fmt.Errorf("user nicht gefunden")
	}

	user := users[0]

	if user.Role == "student" {
		var students []profileRow
		err = database.SupabaseClient.DB.From("students").Select("id").Eq("user_id", user.ID).Execute(&students)
		if err == nil && len(students) > 0 {
			studentID = &students[0].ID
		}
	} else if user.Role == "company" {
		var companies []profileRow
		err = database.SupabaseClient.DB.From("companies").Select("id").Eq("user_id", user.ID).Execute(&companies)
		if err == nil && len(companies) > 0 {
			companyID = &companies[0].ID
		}
	}

	return user.Role, studentID, companyID, nil
}
