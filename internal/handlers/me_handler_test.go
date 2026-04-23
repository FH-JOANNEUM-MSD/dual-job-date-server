package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMyIDHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	rr := httptest.NewRecorder()

	GetMyIDHandler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}
