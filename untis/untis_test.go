package untis

import (
	"os"
	"testing"
	"time"
)

func prepareLogin(t *testing.T) *Untis {
	u, err := NewUntisAPI("bk-ahaus")
	if err != nil {
		panic(err)
	}
	err = u.Login(os.Getenv("UNTIS_USER"), os.Getenv("UNTIS_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func TestUntis_Login(t *testing.T) {
	untis := prepareLogin(t)
	t.Cleanup(func() {
		_ = untis.Logout()
	})
	if untis.token == "" {
		t.Error("Token is empty")
	}
}

func TestUntis_Details(t *testing.T) {
	untis := prepareLogin(t)
	t.Cleanup(func() {
		_ = untis.Logout()
	})
	startDate := time.Now().Add(-(time.Hour + time.Minute*30)) //time.Date(2022, 9, 7, 8, 0, 0, 0, time.Local)
	endDate := time.Now().Add(8 * time.Hour)                   // startDate.Add(time.Hour + time.Minute*30)
	entries, err := untis.Details(startDate, endDate)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected some entries")
	}

}
