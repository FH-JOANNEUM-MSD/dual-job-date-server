package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
	"errors" // 🟢 WICHTIG: Das errors Paket importieren!
	"strconv"
)

// 🟢 HIER IST DIE FEHLENDE DEFINITION:
var ErrCompanyNotFound = errors.New("company nicht gefunden")

// 🟢 NEUER NAME: GetSingleCompanyByID
func GetSingleCompanyByID(id int) (*models.Company, error) {
	var companies []models.Company

	// Wir machen für Supabase wieder einen String draus
	idStr := strconv.Itoa(id)

	err := database.SupabaseClient.DB.From("companies").Select("*").Eq("id", idStr).Execute(&companies)

	if err != nil {
		return nil, err
	}

	if len(companies) == 0 {
		return nil, ErrCompanyNotFound // 🟢 Jetzt kennt Go diesen Error!
	}

	return &companies[0], nil
}
