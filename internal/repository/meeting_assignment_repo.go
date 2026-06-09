package repository

import (
	"dual-job-date-server/internal/database"
	"dual-job-date-server/internal/models"
	"sort"
	"strconv"
	"strings"
)

type AssignMeetingsOptions struct {
	EventID                  int
	SlotIDs                  []int
	StudentIDs               []int
	DryRun                   bool
	IncludeInactiveCompanies bool
	ReplaceExisting          bool
}

type AssignedMeeting struct {
	SlotID         int    `json:"slot_id"`
	StudentID      int    `json:"student_id"`
	CompanyID      int    `json:"company_id"`
	PreferenceType string `json:"preference_type"`
}

type UnassignedCompanySlot struct {
	SlotID    int    `json:"slot_id"`
	CompanyID int    `json:"company_id"`
	Reason    string `json:"reason"`
}

type MeetingAssignmentSummary struct {
	TotalCompanySlots       int `json:"total_company_slots"`
	ExistingMeetingsKept    int `json:"existing_meetings_kept"`
	ExistingMeetingsDeleted int `json:"existing_meetings_deleted"`
	GeneratedMeetings       int `json:"generated_meetings"`
	InsertedMeetings        int `json:"inserted_meetings"`
	AssignedLike            int `json:"assigned_like"`
	AssignedNeutral         int `json:"assigned_neutral"`
	AssignedDislike         int `json:"assigned_dislike"`
	DislikeAvoidedSlots     int `json:"dislike_avoided_slots"`
	UnassignedCompanySlots  int `json:"unassigned_company_slots"`
}

type MeetingAssignmentResult struct {
	DryRun                 bool                     `json:"dry_run"`
	PlannedMeetings        []AssignedMeeting        `json:"planned_meetings"`
	UnassignedCompanySlots []UnassignedCompanySlot  `json:"unassigned_company_slots"`
	InsertedMeetings       int                      `json:"inserted_meetings"`
	Summary                MeetingAssignmentSummary `json:"summary"`
}

func AssignMeetingsByPreferences(opts AssignMeetingsOptions) (MeetingAssignmentResult, error) {
	students, err := getStudentsForAssignment(opts.StudentIDs)
	if err != nil {
		return MeetingAssignmentResult{}, err
	}

	companies, err := getCompaniesForAssignment(opts.IncludeInactiveCompanies)
	if err != nil {
		return MeetingAssignmentResult{}, err
	}

	slots, err := getSlotsForAssignment(opts.EventID, opts.SlotIDs)
	if err != nil {
		return MeetingAssignmentResult{}, err
	}

	preferences, err := getAllPreferences()
	if err != nil {
		return MeetingAssignmentResult{}, err
	}

	existingMeetings, err := getMeetingsForEvent(opts.EventID)
	if err != nil {
		return MeetingAssignmentResult{}, err
	}

	deletedMeetingsCount := 0
	if opts.ReplaceExisting {
		deletedMeetingsCount = len(existingMeetings)
		if !opts.DryRun && deletedMeetingsCount > 0 {
			if err := deleteMeetingsForEvent(opts.EventID); err != nil {
				return MeetingAssignmentResult{}, err
			}
		}
		// Re-assignment starts from scratch when replacement is requested.
		existingMeetings = []models.Meeting{}
	}

	result := assignMeetingsCore(students, companies, slots, preferences, existingMeetings, opts, deletedMeetingsCount)

	if !opts.DryRun && len(result.PlannedMeetings) > 0 {
		if err := insertAssignedMeetings(result.PlannedMeetings, opts.EventID); err != nil {
			return MeetingAssignmentResult{}, err
		}
		result.InsertedMeetings = len(result.PlannedMeetings)
		result.Summary.InsertedMeetings = len(result.PlannedMeetings)
	}

	return result, nil
}

// assignMeetingsCore is the pure greedy matcher: slices in, planned result out, no DB I/O.
func assignMeetingsCore(
	students []models.Student,
	companies []models.Company,
	slots []models.Slot,
	preferences []models.Preference,
	existingMeetings []models.Meeting,
	opts AssignMeetingsOptions,
	deletedMeetingsCount int,
) MeetingAssignmentResult {
	sort.Slice(students, func(i, j int) bool { return students[i].ID < students[j].ID })
	sort.Slice(companies, func(i, j int) bool { return companies[i].ID < companies[j].ID })
	sort.Slice(slots, func(i, j int) bool { return slots[i].ID < slots[j].ID })

	prefMap := buildPreferenceMap(preferences)

	studentLoad := make(map[int]int, len(students))
	companySlotFilled := make(map[string]bool)
	studentSlotTaken := make(map[string]bool)
	studentCompanySeen := make(map[string]bool)

	for _, meeting := range existingMeetings {
		studentLoad[meeting.StudentID]++
		companySlotFilled[companySlotKey(meeting.CompanyID, meeting.SlotID)] = true
		studentSlotTaken[studentSlotKey(meeting.StudentID, meeting.SlotID)] = true
		studentCompanySeen[studentCompanyKey(meeting.StudentID, meeting.CompanyID)] = true
	}

	result := MeetingAssignmentResult{
		DryRun:                 opts.DryRun,
		PlannedMeetings:        []AssignedMeeting{},
		UnassignedCompanySlots: []UnassignedCompanySlot{},
		Summary: MeetingAssignmentSummary{
			TotalCompanySlots:       len(companies) * len(slots),
			ExistingMeetingsDeleted: deletedMeetingsCount,
		},
	}

	for _, slot := range slots {
		for _, company := range companies {
			companySlot := companySlotKey(company.ID, slot.ID)
			if companySlotFilled[companySlot] {
				result.Summary.ExistingMeetingsKept++
				continue
			}

			likeCandidates := make([]int, 0)
			neutralCandidates := make([]int, 0)
			dislikeCandidates := make([]int, 0)

			for _, student := range students {
				if studentSlotTaken[studentSlotKey(student.ID, slot.ID)] {
					continue
				}
				if studentCompanySeen[studentCompanyKey(student.ID, company.ID)] {
					continue
				}

				switch getPreference(prefMap, student.ID, company.ID) {
				case "like":
					likeCandidates = append(likeCandidates, student.ID)
				case "dislike":
					dislikeCandidates = append(dislikeCandidates, student.ID)
				default:
					neutralCandidates = append(neutralCandidates, student.ID)
				}
			}

			if len(likeCandidates) == 0 && len(neutralCandidates) == 0 && len(dislikeCandidates) == 0 {
				result.UnassignedCompanySlots = append(result.UnassignedCompanySlots, UnassignedCompanySlot{
					SlotID:    slot.ID,
					CompanyID: company.ID,
					Reason:    "kein verfuegbarer student (slot bereits belegt oder bereits meeting mit firma)",
				})
				continue
			}

			if len(dislikeCandidates) > 0 && (len(likeCandidates) > 0 || len(neutralCandidates) > 0) {
				result.Summary.DislikeAvoidedSlots++
			}

			selectedPreference := "dislike"
			selectedCandidates := dislikeCandidates

			if len(likeCandidates) > 0 {
				selectedPreference = "like"
				selectedCandidates = likeCandidates
			} else if len(neutralCandidates) > 0 {
				selectedPreference = "neutral"
				selectedCandidates = neutralCandidates
			}

			studentID := pickStudentWithLowestLoad(selectedCandidates, studentLoad)
			meeting := AssignedMeeting{
				SlotID:         slot.ID,
				StudentID:      studentID,
				CompanyID:      company.ID,
				PreferenceType: selectedPreference,
			}

			result.PlannedMeetings = append(result.PlannedMeetings, meeting)
			result.Summary.GeneratedMeetings++

			switch selectedPreference {
			case "like":
				result.Summary.AssignedLike++
			case "neutral":
				result.Summary.AssignedNeutral++
			case "dislike":
				result.Summary.AssignedDislike++
			}

			companySlotFilled[companySlot] = true
			studentSlotTaken[studentSlotKey(studentID, slot.ID)] = true
			studentCompanySeen[studentCompanyKey(studentID, company.ID)] = true
			studentLoad[studentID]++
		}
	}

	result.Summary.UnassignedCompanySlots = len(result.UnassignedCompanySlots)
	return result
}

func getStudentsForAssignment(studentIDs []int) ([]models.Student, error) {
	var students []models.Student

	query := database.SupabaseClient.DB.From("students").Select("*")
	if len(studentIDs) > 0 {
		if err := query.In("id", intsToStrings(studentIDs)).Execute(&students); err != nil {
			return nil, err
		}
		return students, nil
	}

	if err := query.Execute(&students); err != nil {
		return nil, err
	}
	return students, nil
}

func getSlotsForAssignment(eventID int, slotIDs []int) ([]models.Slot, error) {
	var slots []models.Slot

	query := database.SupabaseClient.DB.From("slots").Select("*").Eq("event_id", strconv.Itoa(eventID))
	if len(slotIDs) > 0 {
		if err := query.In("id", intsToStrings(slotIDs)).Execute(&slots); err != nil {
			return nil, err
		}
		return slots, nil
	}

	if err := query.Execute(&slots); err != nil {
		return nil, err
	}
	return slots, nil
}

func getMeetingsForEvent(eventID int) ([]models.Meeting, error) {
	var meetings []models.Meeting
	err := database.SupabaseClient.DB.
		From("meetings").
		Select("*").
		Eq("event_id", strconv.Itoa(eventID)).
		Execute(&meetings)
	if err != nil {
		return nil, err
	}
	return meetings, nil
}

func deleteMeetingsForEvent(eventID int) error {
	var deleted interface{}
	return database.SupabaseClient.DB.
		From("meetings").
		Delete().
		Eq("event_id", strconv.Itoa(eventID)).
		Execute(&deleted)
}

func getCompaniesForAssignment(includeInactive bool) ([]models.Company, error) {
	var companies []models.Company

	query := database.SupabaseClient.DB.From("companies").Select("*")
	if includeInactive {
		err := query.Execute(&companies)
		if err != nil {
			return nil, err
		}
		return companies, nil
	}

	err := query.Eq("active", "true").Execute(&companies)
	if err != nil {
		return nil, err
	}

	return companies, nil
}

func getAllPreferences() ([]models.Preference, error) {
	var preferences []models.Preference

	err := database.SupabaseClient.DB.
		From("preferences").
		Select("*").
		Execute(&preferences)
	if err != nil {
		return nil, err
	}

	return preferences, nil
}

func insertAssignedMeetings(planned []AssignedMeeting, eventID int) error {
	rows := make([]map[string]interface{}, 0, len(planned))
	for _, meeting := range planned {
		rows = append(rows, map[string]interface{}{
			"slot_id":    meeting.SlotID,
			"student_id": meeting.StudentID,
			"company_id": meeting.CompanyID,
			"event_id":   eventID,
		})
	}

	var inserted interface{}
	return database.SupabaseClient.DB.
		From("meetings").
		Insert(rows).
		Execute(&inserted)
}

func intsToStrings(ids []int) []string {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, strconv.Itoa(id))
	}
	return out
}

func buildPreferenceMap(preferences []models.Preference) map[int]map[int]string {
	result := make(map[int]map[int]string)
	for _, preference := range preferences {
		if _, ok := result[preference.StudentID]; !ok {
			result[preference.StudentID] = make(map[int]string)
		}
		result[preference.StudentID][preference.CompanyID] = strings.ToLower(preference.PreferenceType)
	}
	return result
}

func getPreference(prefMap map[int]map[int]string, studentID, companyID int) string {
	if prefByCompany, ok := prefMap[studentID]; ok {
		if preference, ok := prefByCompany[companyID]; ok {
			return preference
		}
	}
	return "neutral"
}

func pickStudentWithLowestLoad(candidates []int, studentLoad map[int]int) int {
	selected := candidates[0]
	selectedLoad := studentLoad[selected]

	for _, candidate := range candidates[1:] {
		load := studentLoad[candidate]
		if load < selectedLoad || (load == selectedLoad && candidate < selected) {
			selected = candidate
			selectedLoad = load
		}
	}

	return selected
}

func companySlotKey(companyID, slotID int) string {
	return strconv.Itoa(companyID) + ":" + strconv.Itoa(slotID)
}

func studentSlotKey(studentID, slotID int) string {
	return strconv.Itoa(studentID) + ":" + strconv.Itoa(slotID)
}

func studentCompanyKey(studentID, companyID int) string {
	return strconv.Itoa(studentID) + ":" + strconv.Itoa(companyID)
}
