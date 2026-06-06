package models

// CreateMeetingInput parst POST /api/meetings.
type CreateMeetingInput struct {
	SlotID    int `json:"slot_id"`
	StudentID int `json:"student_id"`
	CompanyID int `json:"company_id"`
}
