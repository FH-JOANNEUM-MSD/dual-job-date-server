package repository

import (
	"fmt"
	"strings"

	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"

	"github.com/google/uuid"
)

func BulkImportStudentsAndCompanies(input models.BulkImportRequest) (models.BulkImportResult, error) {
	result := models.BulkImportResult{}

	for _, student := range input.Students {
		if err := importStudent(student); err != nil {
			return result, err
		}
		result.StudentsCreated++
	}

	for _, company := range input.Companies {
		if err := importCompany(company); err != nil {
			return result, err
		}
		result.CompaniesCreated++
	}

	return result, nil
}

func importStudent(input models.BulkImportStudent) error {
	if strings.TrimSpace(input.Email) == "" {
		return fmt.Errorf("student email darf nicht leer sein")
	}
	if strings.TrimSpace(input.FirstName) == "" || strings.TrimSpace(input.LastName) == "" {
		return fmt.Errorf("student name darf nicht leer sein")
	}
	if strings.TrimSpace(input.StudyProgram) == "" {
		return fmt.Errorf("study_program darf nicht leer sein")
	}
	studentID, err := getNextNumericID("students")
	if err != nil {
		return err
	}

	authUserUUID := uuid.New().String()
	internalUUID := uuid.New().String()
	userInsert := map[string]interface{}{
		"id":         internalUUID,
		"user_id":    authUserUUID,
		"role":       "student",
		"first_name": input.FirstName,
		"last_name":  input.LastName,
	}
	if err := database.SupabaseClient.DB.From("users").Insert(userInsert).Execute(nil); err != nil {
		return err
	}

	studentInsert := map[string]interface{}{
		"id":            studentID,
		"user_id":       internalUUID,
		"study_program": input.StudyProgram,
	}
	if input.Semester > 0 {
		studentInsert["semester"] = input.Semester
	}
	return database.SupabaseClient.DB.From("students").Insert(studentInsert).Execute(nil)
}

func importCompany(input models.BulkImportCompany) error {
	if strings.TrimSpace(input.Email) == "" {
		return fmt.Errorf("company email darf nicht leer sein")
	}
	if strings.TrimSpace(input.FirstName) == "" || strings.TrimSpace(input.LastName) == "" {
		return fmt.Errorf("company name der ansprechperson darf nicht leer sein")
	}
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("company name darf nicht leer sein")
	}
	companyID, err := getNextNumericID("companies")
	if err != nil {
		return err
	}

	authUserUUID := uuid.New().String()
	internalUUID := uuid.New().String()
	userInsert := map[string]interface{}{
		"id":         internalUUID,
		"user_id":    authUserUUID,
		"role":       "company",
		"first_name": input.FirstName,
		"last_name":  input.LastName,
	}
	if err := database.SupabaseClient.DB.From("users").Insert(userInsert).Execute(nil); err != nil {
		return err
	}

	companyInsert := map[string]interface{}{
		"id":                companyID,
		"user_id":           internalUUID,
		"name":              input.Name,
		"description":       input.Description,
		"short_description": input.ShortDescription,
		"website":           input.Website,
		"logo_url":          input.LogoURL,
		"image_urls":        input.ImageURLs,
		"active":            input.Active,
	}
	return database.SupabaseClient.DB.From("companies").Insert(companyInsert).Execute(nil)
}

func getNextNumericID(table string) (int, error) {
	var rows []struct {
		ID int `json:"id"`
	}
	if err := database.SupabaseClient.DB.From(table).Select("id").Execute(&rows); err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 1, nil
	}
	maxID := rows[0].ID
	for _, row := range rows[1:] {
		if row.ID > maxID {
			maxID = row.ID
		}
	}
	return maxID + 1, nil
}
