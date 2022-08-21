package types

import "time"

const (
	LaunchpadStatusActive = "active"
)

type Launch struct {
	ID        string    `json:"id"`
	DateUTC   time.Time `json:"date_utc"`
	DateLocal time.Time `json:"date_local"`
	Launchpad string    `json:"launchpad"`
}

type Launchpad struct {
	ID       string         `json:"id"`
	FullName string         `json:"full_name"`
	Location *time.Location `json:"location"`
	Status   string         `json:"status"`
}

type LaunchpadFirstDestination struct {
	LaunchpadID   string     `json:"launchpad_id"`
	DestinationID string     `json:"destination_id"`
	LocalYear     int        `json:"local_year"`
	LocalMonth    time.Month `json:"local_month"`
	LocalDay      int        `json:"local_day"`
}
