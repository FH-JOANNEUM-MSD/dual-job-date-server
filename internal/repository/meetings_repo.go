package repository

import (
	"strconv"

	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
)

// GetAllMeetingsWithStudentAndSlot holt alle Meetings für den Admin
func GetAllMeetingsWithStudentAndSlot() ([]models.CompanyMeeting, error) {
	var meetings []models.Meeting
	if err := database.SupabaseClient.DB.From("meetings").Select("*").Execute(&meetings); err != nil {
		return nil, err
	}
	return enrichMeetings(meetings)
}

// GetMeetingsByEvent holt die Meetings eines Events (angereichert mit Slot-Zeiten + Studentennamen).
func GetMeetingsByEvent(eventID int) ([]models.CompanyMeeting, error) {
	var meetings []models.Meeting
	if err := database.SupabaseClient.DB.
		From("meetings").
		Select("*").
		Eq("event_id", strconv.Itoa(eventID)).
		Execute(&meetings); err != nil {
		return nil, err
	}
	return enrichMeetings(meetings)
}

// enrichMeetings joins slot times and student names onto raw meetings.
func enrichMeetings(meetings []models.Meeting) ([]models.CompanyMeeting, error) {
	if len(meetings) == 0 {
		return []models.CompanyMeeting{}, nil
	}

	slotIDs := make(map[int]struct{})
	studentIDs := make(map[int]struct{})
	for _, m := range meetings {
		slotIDs[m.SlotID] = struct{}{}
		studentIDs[m.StudentID] = struct{}{}
	}

	var slots []models.Slot
	if err := database.SupabaseClient.DB.From("slots").Select("*").In("id", intSetToStrings(slotIDs)).Execute(&slots); err != nil {
		return nil, err
	}
	slotByID := make(map[int]models.Slot, len(slots))
	for _, s := range slots {
		slotByID[s.ID] = s
	}

	var studentRows []studentProfileLink
	if err := database.SupabaseClient.DB.From("students").Select("id,user_id").In("id", intSetToStrings(studentIDs)).Execute(&studentRows); err != nil {
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
		if err := database.SupabaseClient.DB.From("users").Select("id,first_name,last_name").In("id", userIDs).Execute(&users); err != nil {
			return nil, err
		}
		for _, u := range users {
			nameByUUID[u.ID] = struct{ first, last string }{u.FirstName, u.LastName}
		}
	}

	studentNameByID := make(map[int]struct{ first, last string }, len(studentRows))
	for _, sr := range studentRows {
		studentNameByID[sr.ID] = nameByUUID[sr.UserID]
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
			EventID:          m.EventID,
		})
	}

	return out, nil
}

func intSetToStrings(set map[int]struct{}) []string {
	out := make([]string, 0, len(set))
	for id := range set {
		out = append(out, strconv.Itoa(id))
	}
	return out
}
