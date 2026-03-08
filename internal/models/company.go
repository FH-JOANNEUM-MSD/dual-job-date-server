package models

type Company struct {
    ID          int    `json:"id"`
    UserID      string `json:"user_id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Website     string `json:"website"`
    LogoURL     string `json:"logo_url"`
    Active      bool   `json:"active"`
    LastUpdated string `json:"last_updated"`
}