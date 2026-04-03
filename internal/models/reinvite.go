package models

// InviteRequest wird für den initialen Invite und den Resend genutzt
type ReinviteRequest struct {
	Email string `json:"email" binding:"required"`
	Role  string `json:"role" binding:"required"` // "student" oder "company"
}
