package repository

import (
	"dual-job-date-server/internal/database"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func HardDeleteStudentAndUser(studentID int) error {
	idStr := strconv.Itoa(studentID)
	log.Printf("[DEBUG] Starte Löschvorgang für Student-ID: %s", idStr)

	// 1. SCHRITT: Hol die interne UUID aus der 'students' Tabelle
	var studentResult []map[string]interface{}
	err := database.SupabaseClient.DB.From("students").Select("user_id").Eq("id", idStr).Execute(&studentResult)

	if err != nil {
		log.Printf("[ERROR] DB-Fehler beim Suchen des Studenten: %v", err)
		return err
	}
	if len(studentResult) == 0 {
		log.Printf("[WARN] Student mit ID %s nicht in 'students' Tabelle gefunden", idStr)
		return fmt.Errorf("student nicht gefunden")
	}

	internalUUID := studentResult[0]["user_id"].(string)
	log.Printf("[DEBUG] Interne UUID gefunden: %s", internalUUID)

	// 2. SCHRITT: Hol die ECHTE Auth-UUID aus der 'users' Tabelle
	var userResult []map[string]interface{}
	err = database.SupabaseClient.DB.From("users").Select("user_id").Eq("id", internalUUID).Execute(&userResult)

	if err != nil {
		log.Printf("[ERROR] DB-Fehler beim Suchen der Auth-UUID: %v", err)
		return err
	}
	if len(userResult) == 0 {
		log.Printf("[WARN] Kein User-Eintrag für interne UUID %s gefunden", internalUUID)
		return fmt.Errorf("user-daten nicht gefunden")
	}

	realAuthUUID := userResult[0]["user_id"].(string)
	log.Printf("[DEBUG] ECHTE Supabase-Auth-UUID identifiziert: %s", realAuthUUID)

	// 3. SCHRITT: Tabellen leeren
	log.Println("[DEBUG] Lösche verknüpfte Meetings und Preferences...")
	database.SupabaseClient.DB.From("meetings").Delete().Eq("student_id", idStr).Execute(nil)
	database.SupabaseClient.DB.From("preferences").Delete().Eq("student_id", idStr).Execute(nil)

	log.Println("[DEBUG] Lösche Eintrag aus 'students'...")
	database.SupabaseClient.DB.From("students").Delete().Eq("id", idStr).Execute(nil)

	log.Println("[DEBUG] Lösche Eintrag aus 'users'...")
	database.SupabaseClient.DB.From("users").Delete().Eq("id", internalUUID).Execute(nil)

	// 4. SCHRITT: Den Auth-User bei Supabase löschen
	log.Printf("[DEBUG] Sende DELETE-Request an Supabase Auth-API für UUID: %s", realAuthUUID)

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	deleteURL := fmt.Sprintf("%s/auth/v1/admin/users/%s", supabaseURL, realAuthUUID)

	req, _ := http.NewRequest("DELETE", deleteURL, nil)
	req.Header.Add("Authorization", "Bearer "+supabaseKey)
	req.Header.Add("apikey", supabaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Netzwerkfehler bei Auth-API: %v", err)
		return err
	}
	defer resp.Body.Close()

	log.Printf("[DEBUG] Supabase Auth-API Antwort-Status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Supabase hat das Löschen abgelehnt. Status: %d", resp.StatusCode)
		return fmt.Errorf("auth-delete fehlgeschlagen (Status %d)", resp.StatusCode)
	}

	log.Printf("[SUCCESS] Student %s (Auth: %s) wurde vollständig eliminiert.", idStr, realAuthUUID)
	return nil
}
