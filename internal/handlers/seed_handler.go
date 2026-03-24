package handlers

import (
	"bytes"
	"dual-job-date-server/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Hilfsfunktion: Legt den User im Supabase Auth an und gibt die ECHTE UUID zurück
func createAuthUser(email, password string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY") // MUSS der service_role key sein!

	url := supabaseURL + "/auth/v1/admin/users"
	payload := map[string]interface{}{
		"email":         email,
		"password":      password,
		"email_confirm": true, // Direkt bestätigen, sonst können sie sich nicht einloggen
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+supabaseKey)
	req.Header.Add("apikey", supabaseKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 422 bedeutet meistens: User existiert schon.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("auth creation failed with status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	authID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("could not parse auth ID")
	}

	return authID, nil
}

func SeedDatabase(w http.ResponseWriter, r *http.Request) {
	insertRows := func(table string, rows []map[string]interface{}) error {
		var result interface{}
		if err := database.SupabaseClient.DB.From(table).Insert(rows).Execute(&result); err != nil {
			return fmt.Errorf("failed to insert into %s: %w", table, err)
		}
		return nil
	}

	var users, students, companies, events, slots, preferences, meetings []map[string]interface{}

	// ---------- GAME OF THRONES STUDENTEN (15) ----------
	studentNames := [][]string{
		{"Jon", "Snow", "Night's Watch Defense"},
		{"Daenerys", "Targaryen", "Dragon Management"},
		{"Arya", "Stark", "Assassination Arts"},
		{"Tyrion", "Lannister", "Political Science"},
		{"Sansa", "Stark", "Kingdom Administration"},
		{"Bran", "Stark", "Historical Visioning"},
		{"Samwell", "Tarly", "Maester Studies"},
		{"Robb", "Stark", "Military Strategy"},
		{"Margaery", "Tyrell", "Public Relations"},
		{"Brienne", "of_Tarth", "Honor & Ethics"},
		{"Jorah", "Mormont", "Exile Survival"},
		{"Theon", "Greyjoy", "Maritime Navigation"},
		{"Gendry", "Baratheon", "Advanced Blacksmithing"},
		{"Podrick", "Payne", "Squire Fundamentals"},
		{"Missandei", "of_Naath", "Multilingual Translation"},
	}

	for i, data := range studentNames {
		// 1. E-Mail generieren (z.B. jon.snow@westeros.com)
		email := strings.ToLower(data[0] + "." + data[1] + "@westeros.com")

		// 2. User ECHT in Supabase Auth anlegen
		authUUID, err := createAuthUser(email, "WinterIsComing123!")
		if err != nil {
			http.Error(w, "Auth Fehler bei "+email+": "+err.Error(), http.StatusInternalServerError)
			return
		}

		idUUID := uuid.New().String()
		studentID := i + 1

		users = append(users, map[string]interface{}{
			"id":         idUUID,
			"user_id":    authUUID, // HIER IST DIE MAGIE!
			"role":       "student",
			"first_name": data[0],
			"last_name":  data[1],
		})

		students = append(students, map[string]interface{}{
			"id":            studentID,
			"user_id":       idUUID,
			"study_program": data[2],
			"semester":      (i % 6) + 1,
		})
	}

	// ---------- GAME OF THRONES COMPANIES (8) ----------
	companyData := []struct {
		First, Last, Name, Desc, Web string
		Active                       bool
	}{
		{"Tywin", "Lannister", "Lannister Gold & Loans", "A Lannister always pays his debts.", "https://lannister.gold", true},
		{"Ned", "Stark", "Winterfell Logistics", "Winter is coming.", "https://winterfell.net", true},
		{"Olenna", "Tyrell", "Highgarden Agriculture", "Growing strong.", "https://highgarden.ag", true},
		{"Euron", "Greyjoy", "Iron Fleet Shipping", "What is dead may never die.", "https://ironfleet.sea", true},
		{"Petyr", "Baelish", "Vale Investments", "Chaos is a ladder.", "https://baelish.info", false},
		{"Jeor", "Mormont", "The Wall Security", "Sword in the darkness.", "https://nightswatch.gov", true},
		{"Roose", "Bolton", "Dreadfort Flaying", "Our blades are sharp.", "https://dreadfort.co", false},
		{"Khal", "Drogo", "Dothraki Equine", "Best horses.", "https://dothraki.horse", true},
	}

	for i, data := range companyData {
		email := strings.ToLower(data.First + "." + data.Last + "@" + strings.Split(data.Web, "://")[1])

		authUUID, err := createAuthUser(email, "WinterIsComing123!")
		if err != nil {
			http.Error(w, "Auth Fehler bei "+email+": "+err.Error(), http.StatusInternalServerError)
			return
		}

		idUUID := uuid.New().String()
		companyID := i + 1

		users = append(users, map[string]interface{}{
			"id":         idUUID,
			"user_id":    authUUID, // UND HIER!
			"role":       "company",
			"first_name": data.First,
			"last_name":  data.Last,
		})

		companies = append(companies, map[string]interface{}{
			"id":           companyID,
			"user_id":      idUUID,
			"name":         data.Name,
			"description":  data.Desc,
			"website":      data.Web,
			"logo_url":     fmt.Sprintf("https://logo.clearbit.com/%s", data.Name),
			"active":       data.Active,
			"last_updated": time.Now().Format(time.RFC3339),
		})
	}

	// ---------- EVENTS (3) ----------
	events = append(events,
		map[string]interface{}{"id": 1, "name": "Westeros Career Fair", "location": "King's Landing", "description": "Meet all the high lords.", "event_date": time.Now().AddDate(0, 1, 0).Format("2006-01-02"), "is_active": true},
		map[string]interface{}{"id": 2, "name": "Essos Start-Up Meet", "location": "Meereen", "description": "Innovation across the narrow sea.", "event_date": time.Now().AddDate(0, 2, 0).Format("2006-01-02"), "is_active": true},
		map[string]interface{}{"id": 3, "name": "Winterfell Networking", "location": "The North", "description": "Cold weather, warm contacts.", "event_date": time.Now().AddDate(0, -1, 0).Format("2006-01-02"), "is_active": false},
	)

	// ---------- SLOTS (16) ----------
	for i := 0; i < 16; i++ {
		hour := 9 + (i / 4)
		minute := (i % 4) * 15
		endMinute := minute + 15
		endHour := hour
		if endMinute == 60 {
			endMinute = 0
			endHour++
		}
		slots = append(slots, map[string]interface{}{
			"id":         i + 1,
			"start_time": fmt.Sprintf("%02d:%02d:00", hour, minute),
			"end_time":   fmt.Sprintf("%02d:%02d:00", endHour, endMinute),
		})
	}

	// ---------- PREFERENCES (120) ----------
	prefTypes := []string{"like", "dislike", "neutral"}
	prefID := 1
	for sID := 1; sID <= 15; sID++ {
		for cID := 1; cID <= 8; cID++ {
			pType := prefTypes[(sID+cID*3)%3]
			preferences = append(preferences, map[string]interface{}{
				"id":              prefID,
				"student_id":      sID,
				"company_id":      cID,
				"preference_type": pType,
			})
			prefID++
		}
	}

	// ---------- MEETINGS (40+) ----------
	meetingID := 1
	for cID := 1; cID <= 8; cID++ {
		for slotID := 1; slotID <= 6; slotID++ {
			sID := ((cID * slotID) % 15) + 1
			meetings = append(meetings, map[string]interface{}{
				"id":         meetingID,
				"slot_id":    slotID,
				"student_id": sID,
				"company_id": cID,
			})
			meetingID++
		}
	}

	// ==========================================
	// INSERTS AUSFÜHREN
	// ==========================================
	if err := insertRows("users", users); err != nil {
		http.Error(w, fmt.Sprintf("Users Error: %v", err), http.StatusInternalServerError)
		return
	}
	if err := insertRows("students", students); err != nil {
		http.Error(w, fmt.Sprintf("Students Error: %v", err), http.StatusInternalServerError)
		return
	}
	if err := insertRows("companies", companies); err != nil {
		http.Error(w, fmt.Sprintf("Companies Error: %v", err), http.StatusInternalServerError)
		return
	}
	if err := insertRows("events", events); err != nil {
		http.Error(w, fmt.Sprintf("Events Error: %v", err), http.StatusInternalServerError)
		return
	}
	if err := insertRows("slots", slots); err != nil {
		http.Error(w, fmt.Sprintf("Slots Error: %v", err), http.StatusInternalServerError)
		return
	}
	if err := insertRows("preferences", preferences); err != nil {
		http.Error(w, fmt.Sprintf("Preferences Error: %v", err), http.StatusInternalServerError)
		return
	}
	if err := insertRows("meetings", meetings); err != nil {
		http.Error(w, fmt.Sprintf("Meetings Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Westeros database successfully seeded! Winter is here 🐺🐉❄️"))
}
