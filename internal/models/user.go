package models

type User struct {
    ID        string `json:"id"`
    UserID    string `json:"user_id"`
    Role      string `json:"role"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}