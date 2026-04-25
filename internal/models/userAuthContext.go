package models

type UserAuthResponse struct {
	ID        string `json:"id"`      // 🟢 NEU: Die echte Datenbank-UUID
	UserID    string `json:"user_id"` // Die Supabase-Auth-ID
	Status    string `json:"status"`
	Role      string `json:"role"`
	StudentID *int   `json:"student_id,omitempty"` // Pointer, damit es null/weggelassen wird, wenn leer
	CompanyID *int   `json:"company_id,omitempty"` // Pointer, damit es null/weggelassen wird, wenn leer
}
