package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
	"errors"

	"github.com/google/uuid"
)

// UpdateUserNames aktualisiert Vor- und Nachname eines Users via Supabase.
func UpdateUserNames(userID uuid.UUID, input models.UpdateUserNameInput) error {
	// 1. Wir bauen uns eine leere Map für die Felder, die sich wirklich ändern
	updateData := make(map[string]interface{})

	// 2. Wir prüfen jeden Pointer: Ist er nicht nil, hat der User ihn im JSON mitgeschickt!
	if input.FirstName != nil {
		updateData["first_name"] = *input.FirstName
	}
	if input.LastName != nil {
		updateData["last_name"] = *input.LastName
	}

	// 3. Sicherheitscheck: Wurde überhaupt etwas mitgeschickt?
	if len(updateData) == 0 {
		return errors.New("keine daten zum updaten gefunden")
	}

	// Eine Variable für die Antwort von Supabase (wir nehmen hier eine generische Map,
	// falls du noch kein komplettes models.User Struct hast. Ansonsten geht auch []models.User)
	var updated []map[string]interface{}

	// 4. Ab an Supabase damit!
	// Wichtig: Supabase erwartet bei Eq() einen String, also machen wir userID.String()
	err := database.SupabaseClient.DB.
		From("users").
		Update(updateData).
		Eq("id", userID.String()).
		Execute(&updated)

	if err != nil {
		return err
	}

	// 5. Optionaler Check: Existiert die ID überhaupt?
	// Supabase wirft bei einem UPDATE auf eine nicht-existierende ID keinen Fehler,
	// es gibt einfach ein leeres Array zurück.
	if len(updated) == 0 {
		return errors.New("kein benutzer mit dieser id gefunden")
	}

	return nil
}
