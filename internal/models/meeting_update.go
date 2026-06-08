package models

// UpdateMeetingInput parst PATCH /api/meetings/{id} (nur gesetzte Felder werden geändert).
type UpdateMeetingInput struct {
	SlotID    *int `json:"slot_id,omitempty"`
	StudentID *int `json:"student_id,omitempty"`
	CompanyID *int `json:"company_id,omitempty"`
	EventID   *int `json:"event_id,omitempty"`
}
