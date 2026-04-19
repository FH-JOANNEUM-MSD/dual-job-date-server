package models

// UpdateStudentInput parst PATCH /api/students/{id} (nur gesetzte Felder werden geändert).
type UpdateStudentInput struct {
	StudyProgram *string `json:"study_program,omitempty"`
	Semester     *int    `json:"semester,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
}
