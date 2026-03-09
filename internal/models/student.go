package models

type Student struct {
    ID           int    `json:"id"`
    UserID       string `json:"user_id"`
    StudyProgram string `json:"study_program"`
    Semester     int    `json:"semester"`
    // Wir fügen diese Felder hinzu, damit sie im JSON-Output landen
    FirstName    string `json:"first_name"`
    LastName     string `json:"last_name"`
}