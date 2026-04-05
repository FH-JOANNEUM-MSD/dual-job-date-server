package models

// UpdateUserNameInput enthält nur die Felder, die aktualisiert werden dürfen.
// Pointer (*string) erlauben uns, fehlende Felder (nil) von leeren Strings ("") zu unterscheiden.
type UpdateUserNameInput struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}
