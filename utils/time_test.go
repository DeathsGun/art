package utils

import (
	"testing"
	"time"
)

func TestTimeFunction(t *testing.T) {
	monday := LerpToPreviousMonday(time.Date(2022, 9, 8, 8, 0, 0, 0, time.Local))
	if monday.Weekday() != time.Monday {
		t.Error("Date should have been monday")
	}
	if monday.Day() != 5 {
		t.Error("Date should have been the 05.09.2022")
	}
}
