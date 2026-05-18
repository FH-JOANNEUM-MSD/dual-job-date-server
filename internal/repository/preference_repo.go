package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
)

// GetAllPreferences holt alle Einträge aus der preferences Tabelle über Supabase
func GetAllPreferences() ([]models.Preference, error) {
	var preferences []models.Preference

	// Holt alle Spalten (*) aus der Tabelle "preferences"
	err := database.SupabaseClient.DB.From("preferences").Select("*").Execute(&preferences)
	if err != nil {
		return nil, err
	}

	// Optional, aber Best Practice:
	// Verhindert, dass 'null' im JSON (Frontend) ankommt, wenn die Tabelle leer ist.
	if len(preferences) == 0 {
		return []models.Preference{}, nil
	}

	return preferences, nil
}
