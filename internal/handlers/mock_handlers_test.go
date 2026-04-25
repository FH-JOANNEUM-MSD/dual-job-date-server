package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetActiveCompaniesHandler_mockOK(t *testing.T) {
	t.Skip("requires repository function mocking for successful path")
}

func TestGetActiveCompaniesHandler_mockError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/companies/active?only_unvoted=not-bool", nil)
	rr := httptest.NewRecorder()

	GetActiveCompaniesHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetAllStudentsHandler_mock(t *testing.T) {
	t.Skip("requires repository function mocking to avoid database dependency")
}

func TestAssignMeetingsByPreferencesHandler_mock(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/meetings/assign", strings.NewReader(`{"dry_run":`))
	rr := httptest.NewRecorder()

	AssignMeetingsByPreferencesHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetPreferencesByStudentHandler_mock(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/students/not-a-number/preferences", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-number"})
	rr := httptest.NewRecorder()

	GetPreferencesByStudentHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
