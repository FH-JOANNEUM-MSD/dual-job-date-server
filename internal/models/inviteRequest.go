package models

type InviteRequest struct {
	Email     string `json:"email"`
	Role      string `json:"role"` // "student" oder "company"
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	// Nur für Studenten relevant
	StudyProgram string `json:"study_program,omitempty"`
	Semester     int    `json:"semester,omitempty"`

	// Nur für Firmen relevant (Beispiele)
	CompanyName string `json:"company_name,omitempty"`
}
