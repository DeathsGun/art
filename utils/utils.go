package utils

import (
	"time"
)

func LeapToPreviousMonday(t time.Time) time.Time {
	return ZeroDate(t).
		Add(time.Hour * 24 * time.Duration(time.Monday-t.Weekday()))
}

func LeapToNearestMonday(t time.Time) time.Time {
	if t.Weekday() <= time.Thursday {
		return LeapToPreviousMonday(t)
	}

	modTime := ZeroDate(t)
	for i := 0; i < 3; i++ {
		modTime = modTime.Add(time.Hour * 24)
		if modTime.Weekday() == time.Monday {
			return ZeroDate(modTime)
		}
	}

	return t
}

func ZeroDate(t time.Time) time.Time {
	return t.Add(-time.Second * time.Duration(t.Second())).
		Add(-time.Minute * time.Duration(t.Minute())).
		Add(-time.Hour * time.Duration(t.Hour()))
}

// Contains is a helper function which checks if
// an element is in an array.
// Because go has no generic array contains function in the
// standard lib.
func Contains[T comparable](array []T, val T) bool {
	for _, t := range array {
		if t == val {
			return true
		}
	}
	return false
}
