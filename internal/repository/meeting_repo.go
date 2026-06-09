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

func getActiveEventID() (int, error) {
	event, err := GetActiveEvent()
	if err != nil {
		return 0, err
	}
	if event.ID == 0 {
		return 0, fmt.Errorf("kein aktives event gefunden")
	}
	return event.ID, nil
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
	return enrichMeetings(meetings)
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
	eventID := current.EventID

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
	if input.EventID != nil {
		if *input.EventID <= 0 {
			return models.Meeting{}, fmt.Errorf("event_id muss größer als 0 sein")
		}
		eventID = *input.EventID
	}

	if err := ensureMeetingReferencesExist(studentID, companyID, slotID); err != nil {
		return models.Meeting{}, err
	}
	if !tableRowExists("events", eventID) {
		return models.Meeting{}, fmt.Errorf("event mit id %d nicht gefunden", eventID)
	}
	if err := ensureMeetingScheduleValid(meetingID, eventID, studentID, companyID, slotID); err != nil {
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
	if input.EventID != nil {
		updateData["event_id"] = eventID
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

func ensureMeetingScheduleValid(meetingID, eventID, studentID, companyID, slotID int) error {
	if taken, err := meetingTakenByOther(meetingID, eventID, "company_id", companyID, "slot_id", slotID); err != nil {
		return err
	} else if taken {
		return meetingConflict("firma hat in diesem slot bereits ein meeting")
	}

	if taken, err := meetingTakenByOther(meetingID, eventID, "student_id", studentID, "slot_id", slotID); err != nil {
		return err
	} else if taken {
		return meetingConflict("student hat in diesem slot bereits ein meeting")
	}

	if taken, err := meetingTakenByOther(meetingID, eventID, "student_id", studentID, "company_id", companyID); err != nil {
		return err
	} else if taken {
		return meetingConflict("student hat bereits ein meeting mit dieser firma")
	}

	return nil
}

// CreateMeeting legt ein einzelnes Meeting an (Admin).
func CreateMeeting(input models.CreateMeetingInput) (models.Meeting, error) {
	var err error

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
	eventID := 0
	if input.EventID != nil {
		if *input.EventID <= 0 {
			return models.Meeting{}, fmt.Errorf("event_id muss größer als 0 sein")
		}
		if _, err := GetEventByID(*input.EventID); err != nil {
			return models.Meeting{}, err
		}
		eventID = *input.EventID
	} else {
		eventID, err = getActiveEventID()
		if err != nil {
			return models.Meeting{}, err
		}
	}
	if err := ensureMeetingScheduleValid(0, eventID, input.StudentID, input.CompanyID, input.SlotID); err != nil {
		return models.Meeting{}, err
	}

	insertData := map[string]interface{}{
		"slot_id":    input.SlotID,
		"student_id": input.StudentID,
		"company_id": input.CompanyID,
		"event_id":   eventID,
	}

	var created []models.Meeting
	err = database.SupabaseClient.DB.From("meetings").Insert(insertData).Execute(&created)
	if err != nil {
		return models.Meeting{}, err
	}
	if len(created) == 0 {
		return models.Meeting{}, fmt.Errorf("meeting konnte nicht angelegt werden")
	}
	if created[0].EventID == 0 {
		created[0].EventID = eventID
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

func meetingTakenByOther(excludeMeetingID, eventID int, field1 string, value1 int, field2 string, value2 int) (bool, error) {
	var rows []models.Meeting
	err := database.SupabaseClient.DB.
		From("meetings").
		Select("id").
		Eq("event_id", strconv.Itoa(eventID)).
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
