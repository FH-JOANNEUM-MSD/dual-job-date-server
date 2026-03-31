package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"

	"github.com/google/uuid"
)

func CreateCompanyProfile(authUUID string, req models.InviteRequest) error {
	internalUUID := uuid.New().String()

	// 1. Eintrag in public.users
	userInsert := map[string]interface{}{
		"id":         internalUUID,
		"user_id":    authUUID,
		"role":       "company",
		"first_name": req.FirstName,
		"last_name":  req.LastName,
	}

	err := database.SupabaseClient.DB.From("users").Insert(userInsert).Execute(nil)
	if err != nil {
		return err
	}

	// 2. Eintrag in public.companies
	companyInsert := map[string]interface{}{
		"user_id": internalUUID,
		"name":    req.CompanyName, // Feldname aus deinem Modell
	}

	err = database.SupabaseClient.DB.From("companies").Insert(companyInsert).Execute(nil)
	return err
}
