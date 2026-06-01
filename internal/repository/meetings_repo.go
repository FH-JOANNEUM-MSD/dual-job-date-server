package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
)

// GetAllMeetingsWithStudentAndSlot holt alle Meetings für den Admin
func GetAllMeetingsWithStudentAndSlot() ([]models.CompanyMeeting, error) {
	var meetings []models.Meeting

	// Keine Filterung (.Eq), holt einfach alle
	err := database.SupabaseClient.DB.From("meetings").Select("*").Execute(&meetings)
	if err != nil {
		return nil, err
	}
	if len(meetings) == 0 {
		return []models.CompanyMeeting{}, nil
	}

	slotIDs := make(map[int]struct{})
	studentIDs := make(map[int]struct{})
	for _, m := range meetings {
		slotIDs[m.SlotID] = struct{}{}
		studentIDs[m.StudentID] = struct{}{}
	}

	slotStrs := intSetToStrings(slotIDs)
	studentStrs := intSetToStrings(studentIDs)

	var slots []models.Slot
	err = database.SupabaseClient.DB.From("slots").Select("*").In("id", slotStrs).Execute(&slots)
	if err != nil {
		return nil, err
	}
	slotByID := make(map[int]models.Slot, len(slots))
	for _, s := range slots {
		slotByID[s.ID] = s
	}

	var studentRows []studentProfileLink
	err = database.SupabaseClient.DB.From("students").Select("id,user_id").In("id", studentStrs).Execute(&studentRows)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(studentRows))
	seenUsers := make(map[string]struct{}, len(studentRows))
	for _, sr := range studentRows {
		if sr.UserID == "" {
			continue
		}
		if _, ok := seenUsers[sr.UserID]; ok {
			continue
		}
		seenUsers[sr.UserID] = struct{}{}
		userIDs = append(userIDs, sr.UserID)
	}

	nameByUUID := map[string]struct{ first, last string }{}
	if len(userIDs) > 0 {
		var users []userDisplayRow
		err = database.SupabaseClient.DB.From("users").Select("id,first_name,last_name").In("id", userIDs).Execute(&users)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			nameByUUID[u.ID] = struct{ first, last string }{u.FirstName, u.LastName}
		}
	}

	studentNameByID := make(map[int]struct{ first, last string }, len(studentRows))
	for _, sr := range studentRows {
		n := nameByUUID[sr.UserID]
		studentNameByID[sr.ID] = n
	}

	out := make([]models.CompanyMeeting, 0, len(meetings))
	for _, m := range meetings {
		slot := slotByID[m.SlotID]
		st := studentNameByID[m.StudentID]
		out = append(out, models.CompanyMeeting{
			ID:               m.ID,
			SlotID:           m.SlotID,
			SlotStartTime:    slot.StartTime,
			SlotEndTime:      slot.EndTime,
			StudentID:        m.StudentID,
			StudentFirstName: st.first,
			StudentLastName:  st.last,
			CompanyID:        m.CompanyID,
		})
	}

	return out, nil
}
