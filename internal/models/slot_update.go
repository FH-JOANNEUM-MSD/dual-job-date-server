package models

// CreateSlotInput parst POST /api/slots. Die ID wird von der Datenbank vergeben.
type CreateSlotInput struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	EventID   int    `json:"event_id"`
}

// UpdateSlotInput parst PATCH /api/slots/{id} (nur gesetzte Felder werden geändert).
type UpdateSlotInput struct {
	StartTime *string `json:"start_time,omitempty"`
	EndTime   *string `json:"end_time,omitempty"`
}
