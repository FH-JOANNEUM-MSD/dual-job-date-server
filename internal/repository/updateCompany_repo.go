package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
	"errors"
	"strconv"
	"strings"
	"time"
)

func UpdateCompany(companyID int, input models.UpdateCompanyInput) error {
	// 1. Wir bauen uns eine leere Map für die Felder, die sich wirklich ändern
	updateData := make(map[string]interface{})
	removedImageURLs := make([]string, 0)

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

		company, err := GetCompanyByID(companyID)
		if err != nil {
			return err
		}

		oldSet := splitImageURLsToSet(company.ImageURLs)
		newSet := splitImageURLsToSet(*input.ImageURLs)
		for oldURL := range oldSet {
			if _, stillPresent := newSet[oldURL]; !stillPresent {
				removedImageURLs = append(removedImageURLs, oldURL)
			}
		}
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
	err := database.SupabaseClient.DB.
		From("companies").
		Update(updateData).
		Eq("id", strconv.Itoa(companyID)).
		Execute(&updated)
	if err != nil {
		return err
	}

	// Best effort: entfernte URLs auch im Image-Bucket loeschen (ohne PATCH zu blockieren).
	for _, removedURL := range removedImageURLs {
		if err := DeleteCompanyImageObjectByURL(removedURL); err != nil {
			// no-op: DB-Update ist erfolgreich, Storage-Cleanup kann spaeter separat passieren
			continue
		}
	}
	return nil
}

func splitImageURLsToSet(raw string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, part := range strings.Split(raw, ";") {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		result[value] = struct{}{}
	}
	return result
}
