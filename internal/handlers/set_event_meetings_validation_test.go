package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestSetEventMeetingsHandler_InvalidEventID_ReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/events/not-a-number/meetings", strings.NewReader(`{"meetings":[]}`))
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-number"})
	rr := httptest.NewRecorder()

	SetEventMeetingsHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestSetEventMeetingsHandler_InvalidMeetingFields_ReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/events/12/meetings", strings.NewReader(`{"meetings":[{"slot_id":0,"student_id":5,"company_id":3}]}`))
	req = mux.SetURLVars(req, map[string]string{"id": "12"})
	rr := httptest.NewRecorder()

	SetEventMeetingsHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}
