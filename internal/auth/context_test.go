package auth

import (
	"context"
	"testing"
)

func TestGetUserID(t *testing.T) {
	t.Run("returns user id from context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", "abc-123")

		got := GetUserID(ctx)
		if got != "abc-123" {
			t.Fatalf("expected userID %q, got %q", "abc-123", got)
		}
	})

	t.Run("returns empty string when missing", func(t *testing.T) {
		got := GetUserID(context.Background())
		if got != "" {
			t.Fatalf("expected empty userID, got %q", got)
		}
	})

	t.Run("returns empty string for wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", 123)

		got := GetUserID(ctx)
		if got != "" {
			t.Fatalf("expected empty userID, got %q", got)
		}
	})
}
