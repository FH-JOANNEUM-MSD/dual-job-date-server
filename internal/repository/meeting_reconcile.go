package repository

import (
	"strconv"

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
