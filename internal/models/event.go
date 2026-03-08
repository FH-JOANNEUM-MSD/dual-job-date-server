package models

type Event struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Location    string `json:"location"`
    Description string `json:"description"`
    EventDate   string `json:"event_date"`
    IsActive    bool   `json:"is_active"`
}