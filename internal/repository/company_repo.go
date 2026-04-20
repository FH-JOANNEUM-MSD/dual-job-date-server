package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
	"strconv"
	"time"
)

func GetActiveCompanies() ([]models.Company, error) {
	var companies []models.Company

	// Wir rufen die Tabelle "companies" ab und filtern nach aktiven Firmen
	err := database.SupabaseClient.DB.From("companies").Select("*").Eq("active", "true").Execute(&companies)

	if err != nil {
		return nil, err
	}

	return companies, nil
}

func GetActiveCompaniesForStudentUser(userID string, onlyUnvoted bool) ([]models.Company, error) {
	companies, err := GetActiveCompanies()
	if err != nil {
		return nil, err
	}

	if !onlyUnvoted {
		return companies, nil
	}

	studentID, err := getStudentIDByUserID(userID)
	if err != nil {
		return nil, err
	}

	var preferences []models.Preference
	err = database.SupabaseClient.DB.
		From("preferences").
		Select("company_id").
		Eq("student_id", strconv.Itoa(studentID)).
		Execute(&preferences)
	if err != nil {
		return nil, err
	}

	votedCompanyIDs := make(map[int]struct{}, len(preferences))
	for _, preference := range preferences {
		votedCompanyIDs[preference.CompanyID] = struct{}{}
	}

	filtered := make([]models.Company, 0, len(companies))
	for _, company := range companies {
		if _, voted := votedCompanyIDs[company.ID]; voted {
			continue
		}
		filtered = append(filtered, company)
	}

	return filtered, nil
}

func GetCompanyByID(companyID int) (models.Company, error) {
	var companies []models.Company

	err := database.SupabaseClient.DB.
		From("companies").
		Select("*").
		Eq("id", strconv.Itoa(companyID)).
		Execute(&companies)
	if err != nil {
		return models.Company{}, err
	}

	if len(companies) == 0 {
		return models.Company{}, nil
	}

	return companies[0], nil
}

func UpdateCompanyLogoURL(companyID int, logoURL string) error {
	update := map[string]interface{}{
		"logo_url":     logoURL,
		"last_updated": time.Now().UTC().Format(time.RFC3339),
	}

	var updated []models.Company
	return database.SupabaseClient.DB.
		From("companies").
		Update(update).
		Eq("id", strconv.Itoa(companyID)).
		Execute(&updated)
}
