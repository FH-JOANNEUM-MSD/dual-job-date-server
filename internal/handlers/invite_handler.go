package handlers

import (
	"dual-job-date-server/internal/models"
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"net/http"
)

func InviteUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.InviteRequest
	json.NewDecoder(r.Body).Decode(&req)

	// 1. User im Auth-System einladen (Supabase verschickt die E-Mail)
	// Diese Funktion schreiben wir gleich im Repository
	authUUID, err := repository.InviteAuthUser(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // <-- Jetzt schickt er die ECHTE Meldung an Bruno!
		return
	}

	// 2. Die "Weiche": Je nach Rolle das richtige Repo aufrufen
	switch req.Role {
	case "student":
		err = repository.CreateStudentProfile(authUUID, req)
	case "company":
		err = repository.CreateCompanyProfile(authUUID, req)
	default:
		http.Error(w, "Unbekannte Rolle", http.StatusBadRequest)
		return
	}

	// 3. Erfolg melden!
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Einladung verschickt und Profil in DB angelegt!"))
}
