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

func TestAssignMeetingsHandler_ReplaceExistingWithSlotIDs_ReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/meetings/assign", strings.NewReader(`{"event_id":12,"replace_existing":true,"slot_ids":[101,102]}`))
	rr := httptest.NewRecorder()

	AssignMeetingsByPreferencesHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestAssignMeetingsHandler_ReplaceExistingWithStudentIDs_ReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/meetings/assign", strings.NewReader(`{"event_id":12,"replace_existing":true,"student_ids":[4,7]}`))
	rr := httptest.NewRecorder()

	AssignMeetingsByPreferencesHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

// A dry_run never deletes (deleteMeetingsForEvent only runs when !dry_run), so
// replace_existing combined with a slot/student subset is SAFE for a preview and
// must NOT be rejected by the guard — the web's "auto-generate" preview relies on it.
// Without a live DB the handler panics/errors once it reaches the repository; we
// recover and only assert that the guard itself did not 400 the request.
func TestAssignMeetingsHandler_DryRunReplaceWithSubset_PassesGuard(t *testing.T) {
	defer func() { _ = recover() }()

	req := httptest.NewRequest(http.MethodPost, "/meetings/assign", strings.NewReader(`{"event_id":12,"dry_run":true,"replace_existing":true,"student_ids":[4,7]}`))
	rr := httptest.NewRecorder()

	AssignMeetingsByPreferencesHandler(rr, req)

	if rr.Code == http.StatusBadRequest && strings.Contains(rr.Body.String(), "replace_existing") {
		t.Fatalf("dry_run request was wrongly rejected by the replace_existing guard: %s", rr.Body.String())
	}
}
