package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateSlotHandler_MissingEventID_ReturnsBadRequest(t *testing.T) {
	body := `{"start_time":"09:00:00","end_time":"09:15:00"}`
	req := httptest.NewRequest(http.MethodPost, "/slots", strings.NewReader(body))
	rr := httptest.NewRecorder()

	CreateSlotHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}
