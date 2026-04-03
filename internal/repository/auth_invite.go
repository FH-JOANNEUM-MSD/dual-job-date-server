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

// InviteAuthUser schickt den Magic-Invite-Link inkl. Redirect-URL und gibt die Supabase Auth-UUID zurück
func InviteAuthUser(email string, role string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	// Wichtig: Hier sollte der SERVICE_ROLE_KEY genutzt werden, damit Admin-Invites erlaubt sind
	supabaseKey := os.Getenv("SUPABASE_KEY")

	inviteURL := supabaseURL + "/auth/v1/invite"

	// --- Logik für die Ziel-URL (Redirect) ---
	var redirectTo string
	if role == "student" {
		// Deep Link für die App (Prüfe nochmal die Schreibweise 'setPassowort'!)
		redirectTo = "dualjob://setPassword"
	} else {
		// Platzhalter für das Web-Team (Company)
		// Sobald die "Lazy Dawgs" die URL liefern, hier anpassen:
		redirectTo = "https://dual-job-webportal-placeholder.de/auth/callback"
	}

	// Body für den Request vorbereiten (jetzt als interface-map für verschiedene Typen)
	bodyData := map[string]interface{}{
		"email":      email,
		"redirectTo": redirectTo,
	}
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

	// 1. Den kompletten Body lesen
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
