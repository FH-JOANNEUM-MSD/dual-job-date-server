package models

import (
	"encoding/json"
	"testing"
)

func TestMeetingDeserializesEventID(t *testing.T) {
	var m Meeting
	if err := json.Unmarshal([]byte(`{"id":1,"slot_id":2,"student_id":3,"company_id":4,"event_id":5}`), &m); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if m.EventID != 5 {
		t.Fatalf("Meeting.EventID = %d, want 5", m.EventID)
	}
}

func TestSlotDeserializesEventID(t *testing.T) {
	var s Slot
	if err := json.Unmarshal([]byte(`{"id":1,"start_time":"09:00:00","end_time":"09:15:00","event_id":7}`), &s); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if s.EventID != 7 {
		t.Fatalf("Slot.EventID = %d, want 7", s.EventID)
	}
}

func TestCompanyMeetingSerializesEventID(t *testing.T) {
	out, err := json.Marshal(CompanyMeeting{EventID: 9})
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if !json.Valid(out) || !containsKey(string(out), `"event_id":9`) {
		t.Fatalf("CompanyMeeting JSON missing event_id: %s", out)
	}
}

func containsKey(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || (len(sub) > 0 && indexOf(s, sub) >= 0))
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
