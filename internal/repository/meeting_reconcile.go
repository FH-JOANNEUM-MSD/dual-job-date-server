package repository

import (
	"strconv"

	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
)

func meetingSeatKey(slotID, studentID, companyID int) string {
	return strconv.Itoa(slotID) + ":" + strconv.Itoa(studentID) + ":" + strconv.Itoa(companyID)
}

// diffMeetings computes which existing meetings to delete and which desired ones to insert
// so that the persisted set equals `desired`. Meetings are matched by (slot, student, company).
func diffMeetings(existing []models.Meeting, desired []AssignedMeeting) (toDeleteIDs []int, toInsert []AssignedMeeting) {
	existingByKey := make(map[string]models.Meeting, len(existing))
	for _, m := range existing {
		existingByKey[meetingSeatKey(m.SlotID, m.StudentID, m.CompanyID)] = m
	}

	desiredByKey := make(map[string]AssignedMeeting, len(desired))
	for _, d := range desired {
		desiredByKey[meetingSeatKey(d.SlotID, d.StudentID, d.CompanyID)] = d
	}

	for key, m := range existingByKey {
		if _, ok := desiredByKey[key]; !ok {
			toDeleteIDs = append(toDeleteIDs, m.ID)
		}
	}
	for key, d := range desiredByKey {
		if _, ok := existingByKey[key]; !ok {
			toInsert = append(toInsert, d)
		}
	}
	return toDeleteIDs, toInsert
}

// SetEventMeetings reconciles the persisted meetings of an event to exactly `desired`.
// Delete-first ordering frees seats before inserting, avoiding swap conflicts.
func SetEventMeetings(eventID int, desired []AssignedMeeting) ([]models.CompanyMeeting, error) {
	existing, err := getMeetingsForEvent(eventID)
	if err != nil {
		return nil, err
	}

	toDeleteIDs, toInsert := diffMeetings(existing, desired)

	for _, id := range toDeleteIDs {
		if err := deleteMeetingByID(id); err != nil {
			return nil, err
		}
	}

	if len(toInsert) > 0 {
		if err := insertAssignedMeetings(toInsert, eventID); err != nil {
			return nil, err
		}
	}

	return GetMeetingsByEvent(eventID)
}

func deleteMeetingByID(meetingID int) error {
	var deleted interface{}
	return database.SupabaseClient.DB.
		From("meetings").
		Delete().
		Eq("id", strconv.Itoa(meetingID)).
		Execute(&deleted)
}
