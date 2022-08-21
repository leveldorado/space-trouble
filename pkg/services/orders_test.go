package services

import (
	context "context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/leveldorado/space-trouble/pkg/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func prepareLaunchpad(t *testing.T) (types.Launchpad, *mockLaunchpadRepo) {
	launchpadLocation, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	launchpad := types.Launchpad{
		ID:       uuid.New().String(),
		Location: launchpadLocation,
	}
	lr := &mockLaunchpadRepo{}
	lr.On("Get", mock.Anything, launchpad.ID).Return(launchpad, nil)
	return launchpad, lr
}

func prepareDestinations() ([]types.Destination, *mockDestinationRepo) {
	destinationsN := 9
	var destinations []types.Destination
	for i := 0; i < destinationsN; i++ {
		destinations = append(destinations, types.Destination{
			ID: uuid.New().String(),
		})
	}
	dr := &mockDestinationRepo{}
	dr.On("ListSorted", mock.Anything).Return(destinations, nil)
	return destinations, dr
}

func prepareFirstDestinationRepo(lID, dID string, year, month, day int) *mockLaunchpadFirstDestinationRepo {
	firstDestination := types.LaunchpadFirstDestination{
		LaunchpadID:   lID,
		DestinationID: dID,
		LocalYear:     year,
		LocalMonth:    month,
		LocalDay:      day,
	}
	lfd := &mockLaunchpadFirstDestinationRepo{}
	lfd.On("Get", mock.Anything, lID).Return(firstDestination, nil)
	return lfd
}

func prepareOrder(t *testing.T, lID, dID string, launchDate time.Time) (types.Order, *mockOrderRepo) {
	o := types.Order{
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		Gender:        gofakeit.Gender(),
		LaunchpadID:   lID,
		DestinationID: dID,
		LaunchDate:    launchDate,
	}
	or := &mockOrderRepo{}
	or.On("Insert", mock.Anything, mock.Anything).
		Return(func(ctx context.Context, doc types.Order) error {
			doc.ID = ""
			o.LaunchDate = o.LaunchDate.UTC()
			o.CreatedAt = doc.CreatedAt
			require.Equal(t, o, doc)
			return nil
		})
	return o, or
}

func prepareCompetitorsLaunchesRepo(launchpad string, localDate time.Time) *mockCompetitorLaunchesRepo {
	clr := &mockCompetitorLaunchesRepo{}
	clr.On("ListByDate", mock.Anything, launchpad, localDate).Return([]types.Launch{}, nil)
	return clr
}

func TestOrders_CreateSuccess(t *testing.T) {
	launchpad, lr := prepareLaunchpad(t)
	destinations, dr := prepareDestinations()
	lfr := prepareFirstDestinationRepo(launchpad.ID, destinations[0].ID, 2053, 3, 3)

	userLocation, err := time.LoadLocation("Europe/Oslo")
	require.NoError(t, err)
	launchDate := time.Date(2053, 3, 5, 1, 0, 0, 0, userLocation)

	o, or := prepareOrder(t, launchpad.ID, destinations[1].ID, launchDate)
	clr := prepareCompetitorsLaunchesRepo(launchpad.ID, launchDate.In(launchpad.Location))

	s := NewOrders(or, lr, dr, lfr, clr)

	_, err = s.Create(context.TODO(), o)
	require.NoError(t, err)

	lr.AssertExpectations(t)
	dr.AssertExpectations(t)
	lfr.AssertExpectations(t)
	or.AssertExpectations(t)
	clr.AssertExpectations(t)
}
