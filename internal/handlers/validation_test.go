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

func TestCreateMeetingHandler_Validation(t *testing.T) {
	tests := []struct {
		name string
		body string
		want int
	}{
		{"missing fields", `{}`, http.StatusBadRequest},
		{"zero slot_id", `{"slot_id":0,"student_id":1,"company_id":1}`, http.StatusBadRequest},
		{"zero student_id", `{"slot_id":1,"student_id":0,"company_id":1}`, http.StatusBadRequest},
		{"zero company_id", `{"slot_id":1,"student_id":1,"company_id":0}`, http.StatusBadRequest},
		{"invalid json", `{"slot_id":`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/meetings", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()
			CreateMeetingHandler(rr, req)
			if rr.Code != tt.want {
				t.Fatalf("expected status %d, got %d", tt.want, rr.Code)
			}
		})
	}
}

func TestDeleteMeetingHandler_Validation(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want int
	}{
		{"invalid id", "not-a-number", http.StatusBadRequest},
		{"zero id", "0", http.StatusBadRequest},
		{"negative id", "-1", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/meetings/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rr := httptest.NewRecorder()
			DeleteMeetingHandler(rr, req)
			if rr.Code != tt.want {
				t.Fatalf("expected status %d, got %d", tt.want, rr.Code)
			}
		})
	}
}

func TestUpdateMeetingHandler_Validation(t *testing.T) {
	tests := []struct {
		name string
		id   string
		body string
		want int
	}{
		{"invalid id", "not-a-number", `{"student_id":1}`, http.StatusBadRequest},
		{"zero id", "0", `{"student_id":1}`, http.StatusBadRequest},
		{"empty body fields", "1", `{}`, http.StatusBadRequest},
		{"invalid json", "1", `{"slot_id":`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPatch, "/meetings/"+tt.id, strings.NewReader(tt.body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rr := httptest.NewRecorder()
			UpdateMeetingHandler(rr, req)
			if rr.Code != tt.want {
				t.Fatalf("expected status %d, got %d", tt.want, rr.Code)
			}
		})
	}
}

func TestGetMeetingsByCompanyHandler_CompanyIDValidation(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want int
	}{
		{"non_numeric", "not-a-number", http.StatusBadRequest},
		{"zero", "0", http.StatusBadRequest},
		{"negative", "-3", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/companies/"+tt.id+"/meetings", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rr := httptest.NewRecorder()
			GetMeetingsByCompanyHandler(rr, req)
			if rr.Code != tt.want {
				t.Fatalf("expected status %d, got %d", tt.want, rr.Code)
			}
		})
	}
}
