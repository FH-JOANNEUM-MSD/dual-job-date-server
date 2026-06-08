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

	internalUUID := uuid.New().String()
	userInsert := map[string]interface{}{
		"id":         internalUUID,
		"user_id":    input.Email,
		"role":       "student",
		"first_name": input.FirstName,
		"last_name":  input.LastName,
	}
	if err := database.SupabaseClient.DB.From("users").Insert(userInsert).Execute(nil); err != nil {
		return err
	}

	studentInsert := map[string]interface{}{
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

	internalUUID := uuid.New().String()
	userInsert := map[string]interface{}{
		"id":         internalUUID,
		"user_id":    input.Email,
		"role":       "company",
		"first_name": input.FirstName,
		"last_name":  input.LastName,
	}
	if err := database.SupabaseClient.DB.From("users").Insert(userInsert).Execute(nil); err != nil {
		return err
	}

	companyInsert := map[string]interface{}{
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
