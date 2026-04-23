package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetMeetingsByStudentHandler_InvalidStudentID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/meetings/student/not-a-number", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-number"})
	rr := httptest.NewRecorder()

	GetMeetingsByStudentHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetPreferencesByStudentHandler_InvalidStudentID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/preferences/student/not-a-number", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-number"})
	rr := httptest.NewRecorder()

	GetPreferencesByStudentHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestVoteCompanyHandler_NoSession_ReturnsUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/companies/1/vote", strings.NewReader(`{"vote":"yes"}`))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	VoteCompanyHandler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestVoteCompanyHandler_InvalidCompanyID_ReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/companies/not-a-number/vote", strings.NewReader(`{"vote":"yes"}`))
	req = req.WithContext(context.WithValue(req.Context(), "userID", "user-1"))
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-number"})
	rr := httptest.NewRecorder()

	VoteCompanyHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestVoteCompanyHandler_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/companies/1/vote", strings.NewReader(`{"vote":`))
	req = req.WithContext(context.WithValue(req.Context(), "userID", "user-1"))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	VoteCompanyHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestInviteUserHandler_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/invite", strings.NewReader(`{invalid`))
	rr := httptest.NewRecorder()

	InviteUserHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestAssignMeetingsByPreferencesHandler_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/meetings/assign", strings.NewReader(`{"dry_run":`))
	rr := httptest.NewRecorder()

	AssignMeetingsByPreferencesHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUploadCompanyImageHandler_InvalidCompanyID_ReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/companies/not-a-number/image", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-number"})
	rr := httptest.NewRecorder()

	UploadCompanyImageHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUploadCompanyLogoHandler_InvalidCompanyID_ReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/companies/not-a-number/logo", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-number"})
	rr := httptest.NewRecorder()

	UploadCompanyLogoHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
