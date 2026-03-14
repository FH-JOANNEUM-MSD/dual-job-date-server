package auth

import (
	"context"
)

// GetUserID holt die Supabase User-ID aus dem Request-Context.
func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return ""
	}
	return userID
}
