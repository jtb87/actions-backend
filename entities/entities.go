package entities

import (
	"time"
)

// StoreInterface interface to retrieve entities from
type StoreInterface interface {
	ActionInterface
	Authentication
	ProfileInterface
	CategoryInterface
}

// Authentication interface holds the interface for authentication
type Authentication interface {
	AuthorizeToken(token string) (Profile, error)
}

// ProfileInterface contract defines all actions on the profile entity
type ProfileInterface interface {
	ProfileAuthentication(p *Profile) error
	// CreateProfile(p *Profile) error
	// UpdateProfile(p *Profile) error
}

// ActionInterface contract defines all the possible interactions with the action entity
type ActionInterface interface {
	GetActionByID(id int) (Action, error)
	UpdateAction(a *Action) error
	CreateAction(a *Action) error
	DeleteAction(id int) error
	GetListOfActions(profileID, CategoryID int) ([]Action, error)
}

// Action holds basic information on when an action has occured
type Action struct {
	ID           int        `json:"id"`
	Subject      string     `json:"subject"`
	Description  string     `json:"description"`
	CategoryID   int        `json:"category_id"`
	CategoryName string     `json:"category_name"`
	ActionDate   *time.Time `json:"action_date"`
	PlannedDate  *time.Time `json:"planned_date"`
	UpdatedAt    *time.Time `json:"updated_at"`
	CreatedAt    *time.Time `json:"created_at"`
	ProfileID    int        `json:"profile_id"`
}

// Profile entity
type Profile struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Token     string `json:"-"`
}

// Goal holds are the information that you need to
// finish something except for a completion date obviously
type Goal struct {
	ID          int
	Name        string
	Description string
}
