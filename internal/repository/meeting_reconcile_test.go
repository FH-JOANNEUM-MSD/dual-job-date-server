package repository

import (
	"sort"
	"testing"

	"dual-job-date-server/internal/models"
)

func TestDiffMeetings_InsertDeleteKeep(t *testing.T) {
	existing := []models.Meeting{
		{ID: 1, SlotID: 100, StudentID: 10, CompanyID: 1}, // keep
		{ID: 2, SlotID: 100, StudentID: 11, CompanyID: 2}, // delete (not desired)
	}
	desired := []AssignedMeeting{
		{SlotID: 100, StudentID: 10, CompanyID: 1}, // matches existing #1 -> keep
		{SlotID: 200, StudentID: 12, CompanyID: 3}, // new -> insert
	}

	toDeleteIDs, toInsert := diffMeetings(existing, desired)

	if len(toDeleteIDs) != 1 || toDeleteIDs[0] != 2 {
		t.Fatalf("toDeleteIDs = %v, want [2]", toDeleteIDs)
	}
	if len(toInsert) != 1 || toInsert[0].SlotID != 200 || toInsert[0].StudentID != 12 || toInsert[0].CompanyID != 3 {
		t.Fatalf("toInsert = %+v, want one (200,12,3)", toInsert)
	}
}

func TestDiffMeetings_EmptyDesiredClearsAll(t *testing.T) {
	existing := []models.Meeting{
		{ID: 1, SlotID: 100, StudentID: 10, CompanyID: 1},
		{ID: 2, SlotID: 200, StudentID: 11, CompanyID: 2},
	}

	toDeleteIDs, toInsert := diffMeetings(existing, nil)

	sort.Ints(toDeleteIDs)
	if len(toDeleteIDs) != 2 || toDeleteIDs[0] != 1 || toDeleteIDs[1] != 2 {
		t.Fatalf("toDeleteIDs = %v, want [1 2]", toDeleteIDs)
	}
	if len(toInsert) != 0 {
		t.Fatalf("toInsert = %+v, want empty", toInsert)
	}
}

func TestDiffMeetings_IdenticalIsNoop(t *testing.T) {
	existing := []models.Meeting{{ID: 1, SlotID: 100, StudentID: 10, CompanyID: 1}}
	desired := []AssignedMeeting{{SlotID: 100, StudentID: 10, CompanyID: 1}}

	toDeleteIDs, toInsert := diffMeetings(existing, desired)

	if len(toDeleteIDs) != 0 || len(toInsert) != 0 {
		t.Fatalf("expected no-op, got delete=%v insert=%+v", toDeleteIDs, toInsert)
	}
}
