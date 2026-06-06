package repository

import (
	"errors"
	"fmt"
	"strconv"

	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
)

var (
	ErrMeetingNotFound = errors.New("meeting nicht gefunden")
	ErrMeetingConflict = errors.New("meeting-konflikt")
)

func meetingConflict(msg string) error {
	return fmt.Errorf("%w: %s", ErrMeetingConflict, msg)
}

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

func GetMeetingByID(meetingID int) (models.Meeting, error) {
	idStr := strconv.Itoa(meetingID)
	var rows []models.Meeting
	err := database.SupabaseClient.DB.From("meetings").Select("*").Eq("id", idStr).Execute(&rows)
	if err != nil {
		return models.Meeting{}, err
	}
	if len(rows) == 0 {
		return models.Meeting{}, ErrMeetingNotFound
	}
	return rows[0], nil
}

// UpdateMeeting ändert slot_id, student_id und/oder company_id eines Meetings (Admin).
func UpdateMeeting(meetingID int, input models.UpdateMeetingInput) (models.Meeting, error) {
	current, err := GetMeetingByID(meetingID)
	if err != nil {
		return models.Meeting{}, err
	}

	slotID := current.SlotID
	studentID := current.StudentID
	companyID := current.CompanyID

	if input.SlotID != nil {
		if *input.SlotID <= 0 {
			return models.Meeting{}, fmt.Errorf("slot_id muss größer als 0 sein")
		}
		slotID = *input.SlotID
	}
	if input.StudentID != nil {
		if *input.StudentID <= 0 {
			return models.Meeting{}, fmt.Errorf("student_id muss größer als 0 sein")
		}
		studentID = *input.StudentID
	}
	if input.CompanyID != nil {
		if *input.CompanyID <= 0 {
			return models.Meeting{}, fmt.Errorf("company_id muss größer als 0 sein")
		}
		companyID = *input.CompanyID
	}

	if err := ensureMeetingReferencesExist(studentID, companyID, slotID); err != nil {
		return models.Meeting{}, err
	}
	if err := ensureMeetingScheduleValid(meetingID, studentID, companyID, slotID); err != nil {
		return models.Meeting{}, err
	}

	updateData := make(map[string]interface{})
	if input.SlotID != nil {
		updateData["slot_id"] = slotID
	}
	if input.StudentID != nil {
		updateData["student_id"] = studentID
	}
	if input.CompanyID != nil {
		updateData["company_id"] = companyID
	}

	idStr := strconv.Itoa(meetingID)
	var updated []models.Meeting
	err = database.SupabaseClient.DB.From("meetings").Update(updateData).Eq("id", idStr).Execute(&updated)
	if err != nil {
		return models.Meeting{}, err
	}
	if len(updated) == 0 {
		return models.Meeting{}, ErrMeetingNotFound
	}
	return updated[0], nil
}

func ensureMeetingReferencesExist(studentID, companyID, slotID int) error {
	if !tableRowExists("students", studentID) {
		return fmt.Errorf("student mit id %d nicht gefunden", studentID)
	}
	if !tableRowExists("companies", companyID) {
		return fmt.Errorf("company mit id %d nicht gefunden", companyID)
	}
	if !tableRowExists("slots", slotID) {
		return fmt.Errorf("slot mit id %d nicht gefunden", slotID)
	}
	return nil
}

func tableRowExists(table string, id int) bool {
	var rows []struct {
		ID int `json:"id"`
	}
	err := database.SupabaseClient.DB.From(table).Select("id").Eq("id", strconv.Itoa(id)).Execute(&rows)
	return err == nil && len(rows) > 0
}

func ensureMeetingScheduleValid(meetingID, studentID, companyID, slotID int) error {
	if taken, err := meetingTakenByOther(meetingID, "company_id", companyID, "slot_id", slotID); err != nil {
		return err
	} else if taken {
		return meetingConflict("firma hat in diesem slot bereits ein meeting")
	}

	if taken, err := meetingTakenByOther(meetingID, "student_id", studentID, "slot_id", slotID); err != nil {
		return err
	} else if taken {
		return meetingConflict("student hat in diesem slot bereits ein meeting")
	}

	if taken, err := meetingTakenByOther(meetingID, "student_id", studentID, "company_id", companyID); err != nil {
		return err
	} else if taken {
		return meetingConflict("student hat bereits ein meeting mit dieser firma")
	}

	return nil
}

// CreateMeeting legt ein einzelnes Meeting an (Admin).
func CreateMeeting(input models.CreateMeetingInput) (models.Meeting, error) {
	if input.SlotID <= 0 {
		return models.Meeting{}, fmt.Errorf("slot_id muss größer als 0 sein")
	}
	if input.StudentID <= 0 {
		return models.Meeting{}, fmt.Errorf("student_id muss größer als 0 sein")
	}
	if input.CompanyID <= 0 {
		return models.Meeting{}, fmt.Errorf("company_id muss größer als 0 sein")
	}

	if err := ensureMeetingReferencesExist(input.StudentID, input.CompanyID, input.SlotID); err != nil {
		return models.Meeting{}, err
	}
	if err := ensureMeetingScheduleValid(0, input.StudentID, input.CompanyID, input.SlotID); err != nil {
		return models.Meeting{}, err
	}

	insertData := map[string]interface{}{
		"slot_id":    input.SlotID,
		"student_id": input.StudentID,
		"company_id": input.CompanyID,
	}

	var created []models.Meeting
	err := database.SupabaseClient.DB.From("meetings").Insert(insertData).Execute(&created)
	if err != nil {
		return models.Meeting{}, err
	}
	if len(created) == 0 {
		return models.Meeting{}, fmt.Errorf("meeting konnte nicht angelegt werden")
	}
	return created[0], nil
}

// DeleteMeeting löscht ein einzelnes Meeting (Admin).
func DeleteMeeting(meetingID int) error {
	if _, err := GetMeetingByID(meetingID); err != nil {
		return err
	}

	idStr := strconv.Itoa(meetingID)
	var deleted interface{}
	return database.SupabaseClient.DB.From("meetings").Delete().Eq("id", idStr).Execute(&deleted)
}

func meetingTakenByOther(excludeMeetingID int, field1 string, value1 int, field2 string, value2 int) (bool, error) {
	var rows []models.Meeting
	err := database.SupabaseClient.DB.
		From("meetings").
		Select("id").
		Eq(field1, strconv.Itoa(value1)).
		Eq(field2, strconv.Itoa(value2)).
		Execute(&rows)
	if err != nil {
		return false, err
	}
	for _, row := range rows {
		if row.ID != excludeMeetingID {
			return true, nil
		}
	}
	return false, nil
}
