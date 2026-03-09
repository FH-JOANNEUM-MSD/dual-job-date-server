package repository

import (
    "dual-job-date-server/internal/database"
    "dual-job-date-server/internal/models"
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