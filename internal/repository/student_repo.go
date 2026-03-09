package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
	"fmt"
)

// GetAllStudents holt alle Studenten und verknüpft sie manuell mit den Namen aus der Users-Tabelle.
// Dies umgeht den Library-Bug beim automatischen Join.
func GetAllStudents() ([]models.Student, error) {
	var students []models.Student
	
	// Hilfsstruktur für die Namen aus der Users-Tabelle
	var userData []struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	fmt.Println("--- DEBUG: Schritt 1 - Hole alle Studenten ---")
	// Schritt 1: Basis-Studentendaten laden
	err := database.SupabaseClient.DB.From("students").
		Select("*").
		Execute(&students)
	
	if err != nil {
		fmt.Printf("--- DEBUG FEHLER SCHRITT 1: %v\n", err)
		return nil, err
	}

	fmt.Println("--- DEBUG: Schritt 2 - Hole alle User-Namen ---")
	// Schritt 2: Namen aus der Users-Tabelle laden (Nutzt 'id' als Match-Feld)
	err = database.SupabaseClient.DB.From("users").
		Select("id, first_name, last_name").
		Execute(&userData)
	
	if err != nil {
		// Wenn hier 'unexpected end of JSON input' kommt, blockiert Supabase (RLS) den Zugriff auf 'users'
		fmt.Printf("--- DEBUG FEHLER SCHRITT 2: %v\n", err)
		return nil, fmt.Errorf("Fehler beim Laden der User-Daten: %v", err)
	}

	// Schritt 3: Map erstellen für schnelles Zusammenführen (ID -> Name)
	nameMap := make(map[string]struct{ first, last string })
	for _, u := range userData {
		nameMap[u.ID] = struct{ first, last string }{u.FirstName, u.LastName}
	}

	// Schritt 4: Studenten-Liste mit Namen befüllen
	for i := range students {
		// Wir schauen nach, ob die user_id des Studenten in unserer nameMap existiert
		if val, ok := nameMap[students[i].UserID]; ok {
			students[i].FirstName = val.first
			students[i].LastName = val.last
		}
	}

	fmt.Printf("--- DEBUG: %d Studenten erfolgreich mit Namen verknüpft ---\n", len(students))
	return students, nil
}