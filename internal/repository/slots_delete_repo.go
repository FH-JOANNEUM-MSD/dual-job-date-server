package repository

import (
	"dual-job-date-server/internal/database"
	"strconv"
)

// DeleteSlot löscht einen Zeitslot und alle darin geplanten Meetings
func DeleteSlot(slotID int) error {
	idStr := strconv.Itoa(slotID)

	// 1. Zuerst alle Meetings löschen, die in diesem Slot stattfinden (Foreign Key!)
	var deletedMeetings []interface{}
	err := database.SupabaseClient.DB.From("meetings").Delete().Eq("slot_id", idStr).Execute(&deletedMeetings)
	if err != nil {
		return err // Abbruch, falls hier was schiefgeht
	}

	// 2. Jetzt, wo der Slot frei ist, können wir ihn selbst gefahrlos löschen
	var deletedSlot []interface{}
	return database.SupabaseClient.DB.From("slots").Delete().Eq("id", idStr).Execute(&deletedSlot)
}
