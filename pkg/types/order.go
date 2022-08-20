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
	Birthday      time.Time `json:"birthday"`
	LaunchpadID   string    `json:"launchpad_id"`
	DestinationID string    `json:"destination_id"`
	LaunchDate    time.Time `json:"launch_date"`
}

func (o Order) Validate() error {
	if o.FirstName == "" {
		return errors.New("first_name is required")
	}
	if o.LastName == "" {
		return errors.New("last_name is required")
	}
	if o.Birthday.IsZero() {
		return errors.New("birthday is required")
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
