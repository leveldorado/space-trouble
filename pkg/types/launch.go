package types

import "time"

type Launch struct {
	ID        string    `json:"id"`
	DateUTC   time.Time `json:"date_utc"`
	DateLocal time.Time `json:"date_local"`
	Launchpad string    `json:"launchpad"`
}

type Launchpad struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}
