package models

type Meeting struct {
    ID        int `json:"id"`
    SlotID    int `json:"slot_id"`
    StudentID int `json:"student_id"`
    CompanyID int `json:"company_id"`
}