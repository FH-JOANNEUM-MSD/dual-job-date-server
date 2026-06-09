package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAssignMeetingsHandler_MissingEventID_ReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/meetings/assign", strings.NewReader(`{"dry_run":true}`))
	rr := httptest.NewRecorder()

	AssignMeetingsByPreferencesHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}
