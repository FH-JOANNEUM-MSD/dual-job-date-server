package models

// CreateEventInput parst POST /api/events. Die ID wird von der Datenbank vergeben.
type CreateEventInput struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	EventDate   string `json:"event_date"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// UpdateEventInput parst PATCH /api/events/{id} (nur gesetzte Felder werden geändert).
type UpdateEventInput struct {
	Name        *string `json:"name,omitempty"`
	Location    *string `json:"location,omitempty"`
	Description *string `json:"description,omitempty"`
	EventDate   *string `json:"event_date,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
