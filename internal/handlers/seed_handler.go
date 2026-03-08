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

    // ---------- USERS ----------
    harryID := uuid.New().String()
    hermioneID := uuid.New().String()
    ronID := uuid.New().String()

    ministryUserID := uuid.New().String()
    gringottsUserID := uuid.New().String()

    users := []map[string]interface{}{
        {"id": harryID, "user_id": uuid.New().String(), "role": "student", "first_name": "Harry", "last_name": "Potter"},
        {"id": hermioneID, "user_id": uuid.New().String(), "role": "student", "first_name": "Hermione", "last_name": "Granger"},
        {"id": ronID, "user_id": uuid.New().String(), "role": "student", "first_name": "Ron", "last_name": "Weasley"},
        {"id": ministryUserID, "user_id": uuid.New().String(), "role": "company", "first_name": "Ministry", "last_name": "Recruiting"},
        {"id": gringottsUserID, "user_id": uuid.New().String(), "role": "company", "first_name": "Gringotts", "last_name": "Bank"},
    }
    if err := insertRows("users", users); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // ---------- STUDENTS ----------
    students := []map[string]interface{}{
        {"user_id": harryID, "study_program": "Defense Against Dark Arts", "semester": 5},
        {"user_id": hermioneID, "study_program": "Magical Law", "semester": 5},
        {"user_id": ronID, "study_program": "Wizard Strategy", "semester": 5},
    }
    if err := insertRows("students", students); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // ---------- COMPANIES ----------
    companies := []map[string]interface{}{
        {"user_id": ministryUserID, "name": "Ministry of Magic", "description": "Government of the Wizarding World", "website": "https://magic.gov"},
        {"user_id": gringottsUserID, "name": "Gringotts Bank", "description": "Wizarding bank run by goblins", "website": "https://gringotts.bank"},
    }
    if err := insertRows("companies", companies); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // ---------- EVENTS ----------
    event := []map[string]interface{}{
        {"name": "Hogwarts Career Fair", "location": "Great Hall", "description": "Meet magical companies", "event_date": time.Now().AddDate(0, 1, 0)},
    }
    if err := insertRows("events", event); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // ---------- SLOTS ----------
    slots := []map[string]interface{}{
        {"start_time": "10:00", "end_time": "10:15"},
        {"start_time": "10:15", "end_time": "10:30"},
        {"start_time": "10:30", "end_time": "10:45"},
    }
    if err := insertRows("slots", slots); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte("Hogwarts database successfully seeded ⚡"))
}