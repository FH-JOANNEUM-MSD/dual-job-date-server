package models

type Company struct {
	ID               int      `json:"id"`
	UserID           string   `json:"user_id"`
	Name             string   `json:"name"`
	ShortDescription string   `json:"short_description"`
	Description      string   `json:"description"`
	Website          string   `json:"website"`
	LogoURL          string   `json:"logo_url"`
	ImageURLs        []string `json:"image_urls"`
	Active           bool     `json:"active"`
	LastUpdated      string   `json:"last_updated"`
}
