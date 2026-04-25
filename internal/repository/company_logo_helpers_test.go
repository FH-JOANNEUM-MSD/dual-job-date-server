package repository

import (
	"testing"
)

func TestNormalizeContentType(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{" IMAGE/PNG ; charset=utf-8 ", "image/png"},
		{"image/jpeg", "image/jpeg"},
		{"image/webp;foo=bar", "image/webp"},
		{"", ""},
	}

	for _, tt := range tests {
		got := normalizeContentType(tt.in)
		if got != tt.want {
			t.Fatalf("normalizeContentType(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestIsSupportedLogoContentType(t *testing.T) {
	supported := []string{"image/png", "image/jpeg", "image/jpg", "image/webp"}
	for _, ct := range supported {
		if !isSupportedLogoContentType(ct) {
			t.Fatalf("expected content type %q to be supported", ct)
		}
	}

	unsupported := []string{"text/plain", "application/json", "", "image/gif"}
	for _, ct := range unsupported {
		if isSupportedLogoContentType(ct) {
			t.Fatalf("expected content type %q to be unsupported", ct)
		}
	}
}

func TestSanitizeFilePart(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{" ACME GmbH ", "acme-gmbh"},
		{"Company_42", "company_42"},
		{"  ---Hello---  ", "hello"},
		{"%%%###", ""},
	}

	for _, tt := range tests {
		got := sanitizeFilePart(tt.in)
		if got != tt.want {
			t.Fatalf("sanitizeFilePart(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
