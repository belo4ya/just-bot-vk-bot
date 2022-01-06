package fuapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	host             string = "https://ruz.fa.ru/api"
	searchEndpoint          = host + "/search"
	scheduleEndpoint        = host + "/schedule"
)

type (
	Group struct {
		ID          string `json:"id"`
		Label       string `json:"label"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}
	Lecturer struct {
		Lecturer          string `json:"lecturer"`
		LecturerCustomUID string `json:"lecturerCustomUID"`
		LecturerEmail     string `json:"lecturerEmail"`
		LecturerOID       int    `json:"lecturerOid"`
		LecturerUID       string `json:"lecturerUID"`
		LecturerRank      string `json:"lecturer_rank"`
		LecturerTitle     string `json:"lecturer_title"`
	}
	ScheduleItem struct {
		Auditorium                string     `json:"auditorium"`
		AuditoriumAmount          int        `json:"auditoriumAmount"`
		AuditoriumOID             int        `json:"auditoriumOid"`
		Author                    string     `json:"author"`
		BeginLesson               string     `json:"beginLesson"`
		Building                  string     `json:"building"`
		BuildingGID               int        `json:"buildingGid"`
		BuildingOID               int        `json:"buildingOid"`
		ContentOfLoadOID          int        `json:"contentOfLoadOid"`
		ContentOfLoadUID          string     `json:"contentOfLoadUID"`
		ContentTableOfLessonsName string     `json:"contentTableOfLessonsName"`
		ContentTableOfLessonsOID  int        `json:"contentTableOfLessonsOid"`
		CreatedDate               string     `json:"createddate"`
		Date                      string     `json:"date"`
		DateOfNest                string     `json:"dateOfNest"`
		DayOfWeek                 int        `json:"dayOfWeek"`
		DayOfWeekString           string     `json:"dayOfWeekString"`
		DetailInfo                string     `json:"detailInfo"`
		Discipline                string     `json:"discipline"`
		DisciplineOID             int        `json:"disciplineOid"`
		DisciplineInPlan          int        `json:"disciplineinplan"`
		DisciplineTypeLoad        int        `json:"disciplinetypeload"`
		Duration                  int        `json:"duration"`
		EndLesson                 string     `json:"endLesson"`
		Group                     string     `json:"group"`
		GroupOID                  int        `json:"groupOid"`
		GroupUID                  string     `json:"groupUID"`
		GroupFacultyOID           int        `json:"group_facultyoid"`
		HideInCapacity            int        `json:"hideincapacity"`
		IsBan                     bool       `json:"isBan"`
		KindOfWork                string     `json:"kindOfWork"`
		KindOfWorkComplexity      int        `json:"kindOfWorkComplexity"`
		KindOfWorkOID             int        `json:"kindOfWorkOid"`
		KindOfWorkUID             string     `json:"kindOfWorkUid"`
		Lecturer                  string     `json:"lecturer"`
		LecturerCustomUID         string     `json:"lecturerCustomUID"`
		LecturerEmail             string     `json:"lecturerEmail"`
		LecturerOID               int        `json:"lecturerOid"`
		LecturerUID               string     `json:"lecturerUID"`
		LecturerRank              string     `json:"lecturer_rank"`
		LecturerTitle             string     `json:"lecturer_title"`
		LessonNumberEnd           int        `json:"lessonNumberEnd"`
		LessonNumberStart         int        `json:"lessonNumberStart"`
		LessonOID                 int        `json:"lessonOid"`
		ListOfLecturers           []Lecturer `json:"listOfLecturers"`
		ModifiedDate              string     `json:"modifieddate"`
		Note                      string     `json:"note"`
		NoteDescription           string     `json:"note_description"`
		ParentSchedule            string     `json:"parentschedule"`
		Replaces                  string     `json:"replaces"`
		Stream                    string     `json:"stream"`
		StreamOID                 int        `json:"streamOid"`
		StreamFacultyOID          int        `json:"stream_facultyoid"`
		SubGroup                  string     `json:"subGroup"`
		SubGroupOID               int        `json:"subGroupOid"`
		SubgroupFacultyOID        int        `json:"subgroup_facultyoid"`
		TableOfLessonsName        string     `json:"tableofLessonsName"`
		TableOfLessonsOID         int        `json:"tableofLessonsOid"`
		URL1                      string     `json:"url1"`
		URL1Description           string     `json:"url1_description"`
		URL2                      string     `json:"url2"`
		URL2Description           string     `json:"url2_description"`
	}
	Schedule []ScheduleItem
)

func GetGroup(name string) (*Group, error) {
	p := fmt.Sprintf("?term=%s&type=group", name)
	url := searchEndpoint + p

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var groups []Group
	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return nil, fmt.Errorf("group with name '%s' not found", name)
	}

	return &groups[0], nil
}

func GetGroupSchedule(groupID string, start time.Time, finish time.Time) (*Schedule, error) {
	p := fmt.Sprintf(
		"/group/%s?start=%s&finish=%s",
		groupID,
		start.Format("2006.01.02"),
		finish.Format("2006.01.02"),
	)
	url := scheduleEndpoint + p

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var schedules Schedule
	if err := json.NewDecoder(resp.Body).Decode(&schedules); err != nil {
		return nil, err
	}

	return &schedules, nil
}
