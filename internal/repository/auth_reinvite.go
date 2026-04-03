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

// ReinviteAuthUser sendet einen frischen Link für bestehende oder neue User
func ReinviteAuthUser(email string, role string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY") // SERVICE_ROLE_KEY nutzen!

	inviteURL := supabaseURL + "/auth/v1/invite"

	var redirectTo string
	if role == "student" {
		// Hier der Branch.io / Smart Link für das Deferred Deep Linking
		redirectTo = "https://dualjob.app.link/invite"
	} else {
		// Web-Portal URL für Companies
		redirectTo = "https://portal.dualjob.de/auth/callback"
	}

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
		return "", fmt.Errorf("netzwerkfehler: %v", err)
	}
	defer resp.Body.Close()

	respBodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("[DEBUG] Resend Error. Status: %d, Body: %s", resp.StatusCode, string(respBodyBytes))
		return "", fmt.Errorf("supabase fehler (status %d)", resp.StatusCode)
	}

	var result map[string]interface{}
	json.Unmarshal(respBodyBytes, &result)

	authID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("id nicht gefunden")
	}

	return authID, nil
}
