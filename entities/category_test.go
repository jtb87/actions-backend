package entities

import (
	"testing"
)

func TestCategoryValidate(t *testing.T) {
	c := Category{Name: "haircut"}
	if got := c.Validate(); got != nil {
		t.Errorf("Category.Validate() = %q, want %v", got, nil)
	}
	c.Name = ""
	if got := c.Validate(); got == nil {
		t.Errorf("Category.Validate() = %q, want %v", got, nil)
	}

	c.Name = "name is way too long to be displayed on a page"
	val := "week"
	c.IntervalType = &val
	if got := c.Validate(); got == nil {
		t.Errorf("Category.Validate() = %q, want %v", got, nil)
	}

}

// func TestDaysSinceLastAction(t *testing.T) {
// 	c := Category{Name: "haircut"}
// 	c.calculateDaysSinceLastAction()
// 	if c.DaysSinceLastAction != nil {
// 		t.Errorf("Category.calculateDaysSinceLastAction() = %v, want %v", c.DaysSinceLastAction, nil)
// 	}
// 	dt := time.Date(2019, time.December, 23, 0, 0, 0, 0, time.UTC)
// 	c = Category{Name: "haircut", LastAction: &Action{
// 		ActionDate: &dt,
// 	}}
// 	c.calculateDaysSinceLastAction()
// 	// zero := 1
// 	// if c.DaysSinceLastAction != &zero {
// 	// 	t.Errorf("Category.calculateDaysSinceLastAction() = %v, want %v", c.DaysSinceLastAction, 1)
// 	// }

// }
