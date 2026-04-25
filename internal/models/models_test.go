package models

import (
	"encoding/json"
	"testing"
)

func TestCompany_JSONRoundTrip(t *testing.T) {
	in := Company{
		ID:               7,
		UserID:           "user-123",
		Name:             "Acme Corp",
		ShortDescription: "short",
		Description:      "long",
		Website:          "https://example.com",
		LogoURL:          "https://example.com/logo.png",
		ImageURLs:        "https://example.com/a.png;https://example.com/b.png",
		Active:           true,
		LastUpdated:      "2026-04-23T00:00:00Z",
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var out Company
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if out != in {
		t.Fatalf("round-trip mismatch: got %#v, want %#v", out, in)
	}
}
