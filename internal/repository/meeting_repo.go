package repository

import (
	"strconv"

	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
)

func GetMeetingsByStudent(studentID int) ([]models.Meeting, error) {
	var meetings []models.Meeting

	err := database.SupabaseClient.DB.From("meetings").Select("*").Eq("student_id", strconv.Itoa(studentID)).Execute(&meetings)
	if err != nil {
		return nil, err
	}

	return meetings, nil
}

type studentProfileLink struct {
	ID     int    `json:"id"`
	UserID string `json:"user_id"`
}

type userDisplayRow struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// GetMeetingsByCompanyWithStudentAndSlot returns meetings for a company with slot times and student names (from linked users row).
func GetMeetingsByCompanyWithStudentAndSlot(companyID int) ([]models.CompanyMeeting, error) {
	var meetings []models.Meeting
	err := database.SupabaseClient.DB.From("meetings").Select("*").Eq("company_id", strconv.Itoa(companyID)).Execute(&meetings)
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

func intSetToStrings(set map[int]struct{}) []string {
	out := make([]string, 0, len(set))
	for id := range set {
		out = append(out, strconv.Itoa(id))
	}
	return out
}
