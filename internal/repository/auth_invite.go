package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// InviteAuthUser schickt den Magic-Invite-Link und gibt die Supabase Auth-UUID zurück
func InviteAuthUser(email string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	inviteURL := supabaseURL + "/auth/v1/invite"

	// Body für den Request vorbereiten
	bodyData := map[string]string{"email": email}
	bodyBytes, _ := json.Marshal(bodyData)

	req, _ := http.NewRequest("POST", inviteURL, bytes.NewBuffer(bodyBytes))
	req.Header.Add("Authorization", "Bearer "+supabaseKey)
	req.Header.Add("apikey", supabaseKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("netzwerkfehler beim invite: %v", err)
	}
	defer resp.Body.Close()

	// 1. WICHTIG: Den kompletten Body sofort lesen!
	respBodyBytes, _ := io.ReadAll(resp.Body)

	// 2. Auf Fehler-Statuscodes prüfen
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("[SUPER-DEBUG] Supabase Auth Error. Status: %d, Body: %s", resp.StatusCode, string(respBodyBytes))
		return "", fmt.Errorf("supabase fehler (status %d): %s", resp.StatusCode, string(respBodyBytes))
	}

	// 3. Bei Erfolg JSON parsen
	var result map[string]interface{}
	if err := json.Unmarshal(respBodyBytes, &result); err != nil {
		return "", fmt.Errorf("json parsing fehler: %v", err)
	}

	authID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("konnte auth_id nicht auslesen. Antwort war: %s", string(respBodyBytes))
	}

	return authID, nil
}
