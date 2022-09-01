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

func TestTimetable(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login("schm132074", os.Getenv("UNTIS_PW"))
	if err != nil {
		t.Fatal(err)
	}
	timetable, err := untis.GetTimetable(&TimetableRequest{
		Id:        untis.loginResponse.PersonId,
		Type:      int(untis.loginResponse.PersonType),
		StartDate: 20220821,
		EndDate:   20220901,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", timetable)
}
