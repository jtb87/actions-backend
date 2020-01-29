package entities

import "time"

// calculateDaysInterval calculate the number of days between two dates
// uses if no `to` is given date of today is taken.
func calculateDaysInterval(from, to *time.Time) int {
	if to == nil {
		t := time.Now()
		to = &t
	}
	result := to.Sub(*from).Hours() / 24
	return int(result)
}
