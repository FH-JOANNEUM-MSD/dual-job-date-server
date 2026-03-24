package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"

	"github.com/google/uuid"
)

func CreateCompanyProfile(authUUID string, req models.InviteRequest) error {
	// 1. Unsere interne ID generieren
	internalUUID := uuid.New().String()

	// 2. Den User in public.users eintragen (Ansprechpartner der Firma)
	userInsert := map[string]interface{}{
		"id":      internalUUID,
		"user_id": authUUID,
		"role":    "company",
	}
	err := database.SupabaseClient.DB.From("users").Insert(userInsert).Execute(nil)
	if err != nil {
		return err
	}

	// 3. Die Firma in public.companies eintragen
	companyInsert := map[string]interface{}{
		"user_id": internalUUID,
		"name":    req.CompanyName,
		// Website, Logo, etc. können später über ein Update befüllt werden
	}
	err = database.SupabaseClient.DB.From("companies").Insert(companyInsert).Execute(nil)

	return err
}
