package handlers

import (
	"dual-job-date-server/internal/models"
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"
)

func InviteUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.InviteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("[ERROR] JSON Decoding failed: %v\n", err)
		http.Error(w, "Ungültiges JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("[INFO] Start Invite für Email: %s, Rolle: %s\n", req.Email, req.Role)

	// 1. User im Auth-System einladen
	// JETZT NEU: Wir geben req.Role mit, damit die Funktion weiß, ob App oder Web-Link!
	authUUID, err := repository.InviteAuthUser(req.Email, req.Role)
	if err != nil {
		fmt.Printf("[ERROR] Auth Invite failed for %s: %v\n", req.Email, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("[DEBUG] Auth User erstellt. UUID: %s\n", authUUID)

	// 2. Die "Weiche" (bleibt gleich, da req.Role hier ja schon genutzt wird)
	switch req.Role {
	case "student":
		fmt.Println("[DEBUG] Versuche Student-Profil anzulegen...")
		err = repository.CreateStudentProfile(authUUID, req)
		if err != nil {
			fmt.Printf("[ERROR] CreateStudentProfile failed: %v\n", err)
			http.Error(w, "Student-Profil konnte nicht erstellt werden: "+err.Error(), http.StatusInternalServerError)
			return
		}
	case "company":
		fmt.Println("[DEBUG] Versuche Company-Profil anzulegen...")
		err = repository.CreateCompanyProfile(authUUID, req)
		if err != nil {
			fmt.Printf("[ERROR] CreateCompanyProfile failed: %v\n", err)
			http.Error(w, "Company-Profil konnte nicht erstellt werden: "+err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		fmt.Printf("[WARN] Unbekannte Rolle empfangen: %s\n", req.Role)
		http.Error(w, "Unbekannte Rolle", http.StatusBadRequest)
		return
	}

	fmt.Printf("[SUCCESS] Alles erledigt für %s\n", req.Email)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Einladung verschickt und Profil in DB angelegt!"))
}
