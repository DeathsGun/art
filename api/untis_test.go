package api

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestLogin(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login(os.Getenv("UNTIS_USER"), os.Getenv("UNTIS_PASSWORD"))
	if err != nil {
		panic(err)
	}
	err = untis.Logout()
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindPerson(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login(os.Getenv("UNTIS_USER"), os.Getenv("UNTIS_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := untis.FindPersonId(StudentElement, "Jona", "Schmidt")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", result)
}

func TestSchoolYear(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login(os.Getenv("UNTIS_USER"), os.Getenv("UNTIS_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := untis.GetCurrentSchoolYear()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", result)
}

func TestSchoolYears(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login(os.Getenv("UNTIS_USER"), os.Getenv("UNTIS_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := untis.GetSchoolYears()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", result)
}

func TestHolidays(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login(os.Getenv("UNTIS_USER"), os.Getenv("UNTIS_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := untis.GetHolidays()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", result)
}

func TestDepartments(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login(os.Getenv("UNTIS_USER"), os.Getenv("UNTIS_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := untis.GetTimeGrids()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", result)
}

func TestRooms(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login(os.Getenv("UNTIS_USER"), os.Getenv("UNTIS_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := untis.GetRooms()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", result)
}

func TestClasses(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login(os.Getenv("UNTIS_USER"), os.Getenv("UNTIS_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := untis.GetClassList()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", result)
}

func TestSubjects(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login(os.Getenv("UNTIS_USER"), os.Getenv("UNTIS_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	subjects, err := untis.GetSubjects()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", subjects)
}

func TestStudents(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login(os.Getenv("UNTIS_USER"), os.Getenv("UNTIS_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	students, err := untis.GetStudents()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", students)
}

func TestTimetable(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login(os.Getenv("UNTIS_USER"), os.Getenv("UNTIS_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	options := &TimetableOptions{
		Element: TimetableRequestElement{
			Id:      strconv.Itoa(untis.loginResponse.PersonId), // the value of the expected key
			Type:    int(untis.loginResponse.PersonType),        // The elementType you search
			KeyType: "id",                                       // The key you want to search for
		},
		StartDate:         20220815,
		EndDate:           20220901,
		OnlyBaseTimetable: false,
		TeacherFields:     []string{"id", "name", "longname", "externalkey"},
	}
	result, err := untis.GetTimetable(options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", result)
}
