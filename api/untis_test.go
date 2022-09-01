package api

import (
	"fmt"
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login("schm132074", os.Getenv("UNTIS_PW"))
	if err != nil {
		panic(err)
	}
	err = untis.Logout()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClasses(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login("schm132074", os.Getenv("UNTIS_PW"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := untis.GetClassList(ClassListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", result)
}

func TestTimetable(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login("schm132074", os.Getenv("UNTIS_PW"))
	if err != nil {
		t.Fatal(err)
	}
	options := TimetableOptions{
		Element: TimetableRequestElement{
			Id:      "schm132074", // the value of the expected key
			Type:    5,            // The elementType you search
			KeyType: "name",       // The key you want to search for
		},
		StartDate:     20220815,
		EndDate:       20220901,
		ShowInfo:      true,
		ClassFields:   []string{"id", "longname", "name"},
		RoomFields:    []string{"id", "longname", "name"},
		SubjectFields: []string{"id", "longname", "name"},
		TeacherFields: []string{"id", "longname", "name"},
	}
	result, err := untis.GetTimetable(options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", result)
}
