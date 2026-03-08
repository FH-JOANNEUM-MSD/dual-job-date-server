package models

type Preference struct {
    ID             int    `json:"id"`
    StudentID      int    `json:"student_id"`
    CompanyID      int    `json:"company_id"`
    PreferenceType string `json:"preference_type"`
}