package entities

import (
	"testing"
	"time"
)

func TestCalculateIntervalDays(t *testing.T) {
	start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC)
	result := calculateDaysInterval(&start, &end)
	if result != 1 {
		t.Errorf("calculateDaysInterval() = %v, want %v", result, 1)
	}
	result = calculateDaysInterval(&start, nil)
	if result < 25 {
		t.Errorf("calculateDaysInterval() = %v, should be bigger then %v", result, 25)
	}
}
