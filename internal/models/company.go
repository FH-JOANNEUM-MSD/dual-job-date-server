package models

type Company struct {
	ID               int    `json:"id"`
	UserID           string `json:"user_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	Website          string `json:"website"`
	LogoURL          string `json:"logo_url"`
	ImageURLs        string `json:"image_urls"` // Später evtl. als Slice/Array parsen
	Active           bool   `json:"active"`
	LastUpdated      string `json:"last_updated"`
}
