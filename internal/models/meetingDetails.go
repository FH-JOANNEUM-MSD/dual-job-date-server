package models

// MeetingDetail ist die "schöne" Ansicht für alle Meeting-Routen
type MeetingDetail struct {
	ID               int    `json:"id"`
	SlotID           int    `json:"slot_id"`
	SlotStartTime    string `json:"slot_start_time"`
	SlotEndTime      string `json:"slot_end_time"`
	StudentID        int    `json:"student_id"`
	StudentFirstName string `json:"student_first_name"`
	StudentLastName  string `json:"student_last_name"`
	CompanyID        int    `json:"company_id"`
	EventID          int    `json:"event_id"`
}
