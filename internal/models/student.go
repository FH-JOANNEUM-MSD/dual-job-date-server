package models

type Student struct {
    ID           int    `json:"id"`
    UserID       string `json:"user_id"`
    StudyProgram string `json:"study_program"`
    Semester     int    `json:"semester"`
}