package models

type BulkImportRequest struct {
	Students  []BulkImportStudent `json:"students,omitempty"`
	Companies []BulkImportCompany `json:"companies,omitempty"`
}

type BulkImportStudent struct {
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	StudyProgram string `json:"study_program"`
	Semester     int    `json:"semester,omitempty"`
}

type BulkImportCompany struct {
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	ShortDescription string `json:"short_description,omitempty"`
	Website          string `json:"website,omitempty"`
	LogoURL          string `json:"logo_url,omitempty"`
	ImageURLs        string `json:"image_urls,omitempty"`
	Active           bool   `json:"active,omitempty"`
}

type BulkImportResult struct {
	StudentsCreated  int `json:"students_created"`
	CompaniesCreated int `json:"companies_created"`
}
