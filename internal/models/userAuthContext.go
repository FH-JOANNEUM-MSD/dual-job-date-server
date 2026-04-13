package models

type UserAuthResponse struct {
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	Role      string `json:"role"`
	StudentID *int   `json:"student_id,omitempty"` // Pointer, damit es null/weggelassen wird, wenn leer
	CompanyID *int   `json:"company_id,omitempty"` // Pointer, damit es null/weggelassen wird, wenn leer
}
