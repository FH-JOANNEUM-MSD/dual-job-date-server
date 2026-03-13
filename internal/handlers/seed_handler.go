package handlers

import (
	"dual-job-date-server/internal/database"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func SeedDatabase(w http.ResponseWriter, r *http.Request) {
	// Helper function to insert a slice of maps into a table using bulk insert
	insertRows := func(table string, rows []map[string]interface{}) error {
		var result interface{}
		if err := database.SupabaseClient.DB.From(table).Insert(rows).Execute(&result); err != nil {
			return fmt.Errorf("failed to insert into %s: %w", table, err)
		}
		return nil
	}

	// Slices für alle Tabellen vorbereiten
	var users []map[string]interface{}
	var students []map[string]interface{}
	var companies []map[string]interface{}
	var events []map[string]interface{}
	var slots []map[string]interface{}
	var preferences []map[string]interface{}
	var meetings []map[string]interface{}

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
		{"Brienne", "of Tarth", "Honor & Ethics"},
		{"Jorah", "Mormont", "Exile Survival"},
		{"Theon", "Greyjoy", "Maritime Navigation"},
		{"Gendry", "Baratheon", "Advanced Blacksmithing"},
		{"Podrick", "Payne", "Squire Fundamentals"},
		{"Missandei", "of Naath", "Multilingual Translation"},
	}

	for i, data := range studentNames {
		idUUID := uuid.New().String()
		userUUID := uuid.New().String()
		studentID := i + 1 // Explizite ID für Relationen

		users = append(users, map[string]interface{}{
			"id":         idUUID,
			"user_id":    userUUID,
			"role":       "student",
			"first_name": data[0],
			"last_name":  data[1],
		})

		students = append(students, map[string]interface{}{
			"id":            studentID,
			"user_id":       idUUID,
			"study_program": data[2],
			"semester":      (i % 6) + 1, // Semester 1 bis 6 gemischt
		})
	}

	// ---------- GAME OF THRONES COMPANIES (8) ----------
	companyData := []struct {
		First, Last, Name, Desc, Web string
		Active                       bool
	}{
		{"Tywin", "Lannister", "Lannister Gold & Loans", "A Lannister always pays his debts.", "https://lannister.gold", true},
		{"Ned", "Stark", "Winterfell Logistics", "Winter is coming. We bring the supplies.", "https://winterfell.net", true},
		{"Olenna", "Tyrell", "Highgarden Agriculture", "Growing strong crops for all of Westeros.", "https://highgarden.ag", true},
		{"Euron", "Greyjoy", "Iron Fleet Shipping", "What is dead may never die.", "https://ironfleet.sea", true},
		{"Petyr", "Baelish", "Vale Investments & Info", "Chaos is a ladder.", "https://baelish.info", false}, // Inactive
		{"Jeor", "Mormont", "The Wall Security", "The sword in the darkness.", "https://nightswatch.gov", true},
		{"Roose", "Bolton", "Dreadfort Flaying Services", "Our blades are sharp.", "https://dreadfort.co", false}, // Inactive
		{"Khal", "Drogo", "Dothraki Equine Exports", "Best horses in Essos.", "https://dothraki.horse", true},
	}

	for i, data := range companyData {
		idUUID := uuid.New().String()
		userUUID := uuid.New().String()
		companyID := i + 1 // Explizite ID

		users = append(users, map[string]interface{}{
			"id":         idUUID,
			"user_id":    userUUID,
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
	// Generiert Slots von 09:00 bis 13:00 im 15-Minuten-Takt
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
	// Jeder Student bewertet jede Company (15 * 8 = 120 Einträge)
	prefTypes := []string{"like", "dislike", "neutral"}
	prefID := 1
	for sID := 1; sID <= 15; sID++ {
		for cID := 1; cID <= 8; cID++ {
			// Pseudo-zufällige Verteilung der Preferences, damit es nicht immer gleich aussieht
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
	// Generiere ein paar fixe Meetings für die Firmen
	meetingID := 1
	for cID := 1; cID <= 8; cID++ {
		for slotID := 1; slotID <= 6; slotID++ { // 6 Meetings pro Firma
			sID := ((cID * slotID) % 15) + 1 // Irgendein Student
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
	// INSERTS AUSFÜHREN (Richtige Reihenfolge!)
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
