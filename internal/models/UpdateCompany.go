package models

// UpdateCompanyInput wird nur verwendet, um Update-Anfragen (PATCH/PUT) aus dem Frontend zu parsen.
type UpdateCompanyInput struct {
	Name             *string `json:"name,omitempty"`
	ShortDescription *string `json:"short_description,omitempty"`
	Description      *string `json:"description,omitempty"`
	Website          *string `json:"website,omitempty"`
	LogoURL          *string `json:"logo_url,omitempty"`
	ImageURLs        *string `json:"image_urls,omitempty"`
	Active           *bool   `json:"active,omitempty"`
}
