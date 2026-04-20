package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
	"errors"
	"strconv"
	"time"
)

func UpdateCompany(companyID int, input models.UpdateCompanyInput) error {
	// 1. Wir bauen uns eine leere Map für die Felder, die sich wirklich ändern
	updateData := make(map[string]interface{})

	// 2. Wir prüfen jeden Pointer: Ist er nicht nil, hat der User ihn im JSON mitgeschickt!
	if input.Name != nil {
		updateData["name"] = *input.Name // Den echten Wert (*input) in die Map legen
	}
	if input.ShortDescription != nil {
		updateData["short_description"] = *input.ShortDescription
	}
	if input.Description != nil {
		updateData["description"] = *input.Description
	}
	if input.Website != nil {
		updateData["website"] = *input.Website
	}
	if input.LogoURL != nil {
		updateData["logo_url"] = *input.LogoURL
	}
	if input.ImageURLs != nil {
		updateData["image_urls"] = *input.ImageURLs
	}
	if input.Active != nil {
		updateData["active"] = *input.Active
	}

	// 3. Sicherheitscheck: Wurde überhaupt etwas mitgeschickt?
	if len(updateData) == 0 {
		return errors.New("keine daten zum updaten gefunden")
	}

	// 4. Den Zeitstempel bei jedem Update automatisch erneuern
	updateData["last_updated"] = time.Now().UTC().Format(time.RFC3339)

	var updated []models.Company

	// 5. Ab an Supabase damit!
	return database.SupabaseClient.DB.
		From("companies").
		Update(updateData).
		Eq("id", strconv.Itoa(companyID)).
		Execute(&updated)
}
