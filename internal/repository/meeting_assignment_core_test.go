package repository

import (
	"testing"

	"dual-job-date-server/internal/models"
)

func TestAssignMeetingsCore_FillsSeatsAndPrefersLikes(t *testing.T) {
	students := []models.Student{{ID: 1}, {ID: 2}}
	companies := []models.Company{{ID: 10}, {ID: 20}}
	slots := []models.Slot{{ID: 100}, {ID: 200}}
	// student 1 likes company 10; everything else neutral (absent)
	prefs := []models.Preference{{StudentID: 1, CompanyID: 10, PreferenceType: "like"}}

	res := assignMeetingsCore(students, companies, slots, prefs, nil, AssignMeetingsOptions{}, 0)

	// 2 companies x 2 slots = 4 seats, enough students to fill all (2 students, no double-book per slot/company)
	if res.Summary.TotalCompanySlots != 4 {
		t.Fatalf("TotalCompanySlots = %d, want 4", res.Summary.TotalCompanySlots)
	}
	if len(res.PlannedMeetings) != 4 {
		t.Fatalf("planned meetings = %d, want 4", len(res.PlannedMeetings))
	}
	// The like (student 1 ↔ company 10) must be honored in at least one planned meeting.
	foundLike := false
	for _, m := range res.PlannedMeetings {
		if m.StudentID == 1 && m.CompanyID == 10 && m.PreferenceType == "like" {
			foundLike = true
		}
	}
	if !foundLike {
		t.Fatalf("expected student 1 ↔ company 10 assigned as a like; planned = %+v", res.PlannedMeetings)
	}
	if res.InsertedMeetings != 0 {
		t.Fatalf("core must not report inserts; InsertedMeetings = %d, want 0", res.InsertedMeetings)
	}
}

func TestAssignMeetingsCore_DislikeIsLastResort(t *testing.T) {
	students := []models.Student{{ID: 1}}
	companies := []models.Company{{ID: 10}}
	slots := []models.Slot{{ID: 100}}
	prefs := []models.Preference{{StudentID: 1, CompanyID: 10, PreferenceType: "dislike"}}

	res := assignMeetingsCore(students, companies, slots, prefs, nil, AssignMeetingsOptions{}, 0)

	// Only one possible pairing and it is disliked: it must still be filled (last resort).
	if len(res.PlannedMeetings) != 1 || res.PlannedMeetings[0].PreferenceType != "dislike" {
		t.Fatalf("expected 1 dislike assignment as last resort; got %+v", res.PlannedMeetings)
	}
	if res.Summary.AssignedDislike != 1 {
		t.Fatalf("AssignedDislike = %d, want 1", res.Summary.AssignedDislike)
	}
}

func TestAssignMeetingsCore_RespectsExistingMeetings(t *testing.T) {
	students := []models.Student{{ID: 1}}
	companies := []models.Company{{ID: 10}}
	slots := []models.Slot{{ID: 100}}
	existing := []models.Meeting{{ID: 99, SlotID: 100, StudentID: 1, CompanyID: 10}}

	res := assignMeetingsCore(students, companies, slots, nil, existing, AssignMeetingsOptions{}, 0)

	if len(res.PlannedMeetings) != 0 {
		t.Fatalf("the only seat is already filled; planned should be empty, got %+v", res.PlannedMeetings)
	}
	if res.Summary.ExistingMeetingsKept != 1 {
		t.Fatalf("ExistingMeetingsKept = %d, want 1", res.Summary.ExistingMeetingsKept)
	}
}
