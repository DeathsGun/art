package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLeapToNearestMonday(t *testing.T) {
	t.Parallel()

	input, err := time.Parse("02.01.2006", "22.10.2022")
	if !assert.NoError(t, err) {
		return
	}

	expectedMonday, err := time.Parse("02.01.2006", "24.10.2022")
	if !assert.NoError(t, err) {
		return
	}

	monday := LeapToNearestMonday(input)
	if assert.Equal(t, expectedMonday, monday) {
		return
	}
}

func TestLeapToPreviousMonday(t *testing.T) {
	t.Parallel()

	input, err := time.Parse("02.01.2006", "22.10.2022")
	if !assert.NoError(t, err) {
		return
	}

	expectedMonday, err := time.Parse("02.01.2006", "17.10.2022")
	if !assert.NoError(t, err) {
		return
	}

	monday := LeapToPreviousMonday(input)
	if assert.Equal(t, expectedMonday, monday) {
		return
	}
}
