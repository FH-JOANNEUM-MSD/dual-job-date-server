package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GetAllStudents holt alle Studenten und verknüpft sie mit den User-Daten
func GetAllStudents() ([]models.Student, error) {
	var students []models.Student

	fmt.Println("--- DEBUG: Lade Studenten ---")
	err := database.SupabaseClient.DB.From("students").Select("*").Execute(&students)
	if err != nil {
		fmt.Printf("--- FEHLER STUDENTS: %v\n", err)
		return nil, err
	}

	fmt.Println("--- DEBUG: Lade User per RAW HTTP ---")

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	url := fmt.Sprintf("%s/rest/v1/users?select=id,first_name,last_name", supabaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen des Requests: %v", err)
	}

	req.Header.Add("apikey", supabaseKey)
	req.Header.Add("Authorization", "Bearer "+supabaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP Request fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()

	// Wir lesen die GANZ ROHE Antwort von Supabase
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	fmt.Printf("--- RAW SUPABASE ANTWORT USERS (Status %d) ---\n%s\n---------------------------------------\n", resp.StatusCode, bodyString)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("supabase API Fehler: %s", bodyString)
	}

	var rawUsers []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &rawUsers)
	if err != nil {
		fmt.Printf("--- JSON PARSE FEHLER: %v\n", err)
		return nil, err
	}

	// Mapping bauen
	nameMap := make(map[string]struct{ first, last string })
	for _, u := range rawUsers {
		id, _ := u["id"].(string)
		fName, _ := u["first_name"].(string)
		lName, _ := u["last_name"].(string)

		if id != "" {
			nameMap[id] = struct{ first, last string }{fName, lName}
		}
	}

	// Zusammenführen
	for i := range students {
		if val, ok := nameMap[students[i].UserID]; ok {
			students[i].FirstName = val.first
			students[i].LastName = val.last
		}
	}

	fmt.Printf("--- ERFOLG: %d Studenten und %d User gemappt ---\n", len(students), len(rawUsers))
	return students, nil
}
