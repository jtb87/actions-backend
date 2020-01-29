package entities

import (
	"fmt"
	"time"
)

// CategoryInterface contract defines all the possible interactions with the category entity
type CategoryInterface interface {
	GetListOfCategories(profileID int) ([]Category, error)
	// UpdateCategory(a *Category) error
	CreateCategory(c *Category) error
	DeleteCategory(id int) error
	GetCategory(id int) (Category, error)
}

// Category entity
type Category struct {
	ID                  int        `json:"id"`
	Name                string     `json:"name"`
	IntervalType        *string    `json:"interval_type"`
	Interval            *int       `json:"interval"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DaysSinceLastAction *int       `json:"days_since_last_action"`

	LastAction *Action `json:"-"`
}

// CalcDaysSinceLastAction calculates the DaysSinceLastActionField
func (c *Category) CalcDaysSinceLastAction() {
	if c.LastAction != nil && c.LastAction.ActionDate != nil {
		interval := calculateDaysInterval(c.LastAction.ActionDate, nil)
		c.DaysSinceLastAction = &interval
	}
}

const nameLength = 24

// Validate validates everything
func (c *Category) Validate() (err error) {
	if len(c.Name) > nameLength {
		return fmt.Errorf("Name is too long")
	}
	if len(c.Name) < 1 {
		return fmt.Errorf("Name is too short")
	}

	if c.validIntervalType() != true {
		return fmt.Errorf("%s is not a valid interval_type", *c.IntervalType)
	}
	return
}

func (c *Category) validIntervalType() bool {
	if c.IntervalType == nil {
		return true
	}
	switch *c.IntervalType {
	case
		"hour",
		"day",
		"week",
		"month",
		"year":
		return true
	}
	return false
}

func (c *Category) calculateDaysSinceLastAction() {}
