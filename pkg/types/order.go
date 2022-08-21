package types

import (
	"time"

	"github.com/pkg/errors"
)

type Order struct {
	ID            string    `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Gender        string    `json:"gender"`
	BirthdayYear  int       `json:"birthday_year"`
	BirthdayMonth int       `json:"birthday_month"`
	BirthdayDay   int       `json:"birthday_day"`
	LaunchpadID   string    `json:"launchpad_id"`
	DestinationID string    `json:"destination_id"`
	LaunchDate    time.Time `json:"launch_date"`
	CreatedAt     time.Time `json:"created_at"`
}

func (o Order) Validate() error {
	if o.FirstName == "" {
		return errors.New("first_name is required")
	}
	if o.LastName == "" {
		return errors.New("last_name is required")
	}
	if o.BirthdayYear == 0 {
		return errors.New("birthday year is required")
	}
	if o.BirthdayMonth == 0 {
		return errors.New("birthday month is required")
	}
	if o.BirthdayDay == 0 {
		return errors.New("birthday day is required")
	}
	if o.LaunchpadID == "" {
		return errors.New("launchpad_id is required")
	}
	if o.DestinationID == "" {
		return errors.New("destination_id is required")
	}
	return nil
}

type Destination struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
