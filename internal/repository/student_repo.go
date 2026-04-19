package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
)

// ErrStudentNotFound wird zurückgegeben, wenn keine Zeile in students zur ID existiert.
var ErrStudentNotFound = errors.New("student nicht gefunden")

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

// UpdateStudent wendet partielle Änderungen auf students an und optional Vor-/Nachname in users.
func UpdateStudent(studentID int, input models.UpdateStudentInput) error {
	idStr := strconv.Itoa(studentID)

	var rows []struct {
		UserID string `json:"user_id"`
	}
	err := database.SupabaseClient.DB.From("students").Select("user_id").Eq("id", idStr).Execute(&rows)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return ErrStudentNotFound
	}
	internalUserUUID := rows[0].UserID

	updateData := make(map[string]interface{})
	if input.StudyProgram != nil {
		updateData["study_program"] = *input.StudyProgram
	}
	if input.Semester != nil {
		updateData["semester"] = *input.Semester
	}

	if len(updateData) > 0 {
		var updated []map[string]interface{}
		err = database.SupabaseClient.DB.From("students").Update(updateData).Eq("id", idStr).Execute(&updated)
		if err != nil {
			return err
		}
		if len(updated) == 0 {
			return ErrStudentNotFound
		}
	}

	if input.FirstName != nil || input.LastName != nil {
		uid, err := uuid.Parse(internalUserUUID)
		if err != nil {
			return fmt.Errorf("ungültige user_id für student: %w", err)
		}
		return UpdateUserNames(uid, models.UpdateUserNameInput{
			FirstName: input.FirstName,
			LastName:  input.LastName,
		})
	}

	return nil
}
