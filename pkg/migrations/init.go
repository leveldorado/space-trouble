package migrations

import (
	"context"
	"time"

	"github.com/leveldorado/space-trouble/pkg/repositories"
	"github.com/leveldorado/space-trouble/pkg/types"
	"github.com/pkg/errors"
)

/*
Init initialize things like table creation and populating data.

	in production app better approach to use dedicated tool for it
	or code which triggered not on start but by event so only one replica of app will perform
*/
func Init(
	or *repositories.PostgreSQLOrdersRepo,
	lr *repositories.SpaceXAPILaunchpadsRepo,
	dr *repositories.InMemoryDestinationsRepo,
	fr *repositories.InMemoryLaunchpadFirstDestinationRepo,
) error {
	if err := or.CreateTables(context.TODO()); err != nil {
		return errors.Wrap(err, `failed to create tables`)
	}
	return populateLaunchpadFirstDestinations(lr, dr, fr)
}

/*
launchpad first destination records needed as starting point of calculating destination for a date
*/
func populateLaunchpadFirstDestinations(
	lr *repositories.SpaceXAPILaunchpadsRepo,
	dr *repositories.InMemoryDestinationsRepo,
	fr *repositories.InMemoryLaunchpadFirstDestinationRepo,
) error {
	launchpads, err := lr.List(context.TODO())
	if err != nil {
		return errors.Wrap(err, `failed to list launchpads`)
	}
	destinations, err := dr.ListSorted(context.TODO())
	if err != nil {
		return errors.Wrap(err, `failed to list destinations`)
	}
	var currentDestinationIndex int
	for _, pad := range launchpads {
		padTime := time.Now().In(pad.Location)
		year, month, day := padTime.Date()
		fr.Set(types.LaunchpadFirstDestination{
			LaunchpadID:   pad.ID,
			DestinationID: destinations[0].ID,
			LocalYear:     year,
			LocalMonth:    month,
			LocalDay:      day,
		})
		currentDestinationIndex++
		if currentDestinationIndex >= len(destinations) {
			currentDestinationIndex = 0
		}
	}
	return nil
}
