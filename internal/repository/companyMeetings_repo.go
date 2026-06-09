package repository

import (
	"strconv"

	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
)

// GetMeetingsByCompanyID holt alle Meetings für eine spezifische Firma
func GetMeetingsByCompanyID(companyID int) ([]models.Meeting, error) {
	var meetings []models.Meeting

	// Filtert die Tabelle mit .Eq() nach der company_id
	err := database.SupabaseClient.DB.From("meetings").
		Select("*").
		Eq("company_id", strconv.Itoa(companyID)).
		Execute(&meetings)

	if err != nil {
		return nil, err
	}

	// Leeres Array zurückgeben statt null
	if len(meetings) == 0 {
		return []models.Meeting{}, nil
	}

	for i := range meetings {
		if meetings[i].EventID == 0 {
			activeEventID, err := getActiveEventID()
			if err != nil {
				return nil, err
			}
			meetings[i].EventID = activeEventID
		}
	}

	return meetings, nil
}
