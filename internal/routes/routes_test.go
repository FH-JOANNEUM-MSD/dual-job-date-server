package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthRoot(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Server läuft!") {
		t.Fatalf("expected response to contain %q, got %q", "Server läuft!", body)
	}
}

func TestNewRouter_ProtectedApiRoute_RequiresJWT(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/events/active", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestNewRouter_ResendInviteRoute_IsPublicAndParsesBody(t *testing.T) {
	router := NewRouter()

	// Invalid JSON should fail in handler itself (400),
	// proving this route is not blocked by JWT middleware.
	req := httptest.NewRequest(http.MethodPost, "/api/resend-invite", strings.NewReader("{invalid"))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
