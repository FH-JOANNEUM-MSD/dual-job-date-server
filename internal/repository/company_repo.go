package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
	"strconv"
	"strings"
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

func AddCompanyImageURL(companyID int, imageURL string) error {
	company, err := GetCompanyByID(companyID)
	if err != nil {
		return err
	}
	if company.ID == 0 {
		return nil
	}

	rawEntries := strings.Split(company.ImageURLs, ";")
	imageURLs := make([]string, 0, len(rawEntries)+1)
	seen := make(map[string]struct{}, len(rawEntries)+1)
	for _, raw := range rawEntries {
		existing := strings.TrimSpace(raw)
		if existing == "" {
			continue
		}
		if _, ok := seen[existing]; ok {
			continue
		}
		seen[existing] = struct{}{}
		imageURLs = append(imageURLs, existing)
	}
	candidate := strings.TrimSpace(imageURL)
	if _, ok := seen[candidate]; !ok && candidate != "" {
		imageURLs = append(imageURLs, candidate)
	}

	update := map[string]interface{}{
		"image_urls":   strings.Join(imageURLs, ";"),
		"last_updated": time.Now().UTC().Format(time.RFC3339),
	}
	var updated []models.Company
	return database.SupabaseClient.DB.
		From("companies").
		Update(update).
		Eq("id", strconv.Itoa(companyID)).
		Execute(&updated)
}
