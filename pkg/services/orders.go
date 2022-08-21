package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/leveldorado/space-trouble/pkg/types"
	"github.com/pkg/errors"
)

type orderRepo interface {
	Get(ctx context.Context, id string) (types.Order, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]types.Order, error)
	Insert(ctx context.Context, o types.Order) error
}

type launchpadRepo interface {
	Get(ctx context.Context, id string) (types.Launchpad, error)
}

type destinationRepo interface {
	ListSorted(ctx context.Context) ([]types.Destination, error)
}

type launchpadFirstDestinationRepo interface {
	Get(ctx context.Context, launchpad string) (types.LaunchpadFirstDestination, error)
}

type competitorLaunchesRepo interface {
	CheckLaunches(ctx context.Context, launchpad string, localDate time.Time) (bool, error)
}

type Orders struct {
	orderRepo                     orderRepo
	launchpadRepo                 launchpadRepo
	destinationRepo               destinationRepo
	launchpadFirstDestinationRepo launchpadFirstDestinationRepo
	competitorLaunchesRepo        competitorLaunchesRepo
}

func NewOrders(
	or orderRepo,
	lr launchpadRepo,
	dr destinationRepo,
	lfr launchpadFirstDestinationRepo,
	cr competitorLaunchesRepo,
) *Orders {
	return &Orders{
		orderRepo:                     or,
		launchpadRepo:                 lr,
		destinationRepo:               dr,
		launchpadFirstDestinationRepo: lfr,
		competitorLaunchesRepo:        cr,
	}
}

/*
Create
validate launchpad id and destination id
validate if date is not in past

launchpad destination on date calculates by logic

		we have stored data about first destination from launchpad and date
	    then we calculate difference in days between first launch and requested date
	    and shift destinations by diff days in destination list (required destinations to be sorted)
*/
func (s *Orders) Create(ctx context.Context, o types.Order) (string, error) {
	launchpad, err := s.launchpadRepo.Get(ctx, o.LaunchpadID)
	if errors.As(err, &types.ErrNotFound{}) {
		return "", types.NewErrInvalidData("invalid launchpad id")
	}
	if err != nil {
		return "", errors.Wrapf(err, `failed to get launchpad: id - %s`, o.LaunchpadID)
	}
	if err := s.checkLaunchpadDestination(ctx, launchpad, o); err != nil {
		return "", err
	}
	exists, err := s.competitorLaunchesRepo.CheckLaunches(ctx, o.LaunchpadID, o.LaunchDate.In(launchpad.Location))
	if err != nil {
		return "", errors.Wrapf(err, `failed to list competitor launches by date: launchpad - %s, date - %s`, o.LaunchDate, o.LaunchDate)
	}
	if exists {
		return "", types.ErrFlightImpossible{}
	}
	o.ID = uuid.New().String()
	o.LaunchDate = o.LaunchDate.UTC()
	o.CreatedAt = time.Now().UTC()
	return o.ID, errors.Wrapf(s.orderRepo.Insert(ctx, o), `failed to insert order: o - %+v`, o)
}

func (s *Orders) checkLaunchpadDestination(ctx context.Context, launchpad types.Launchpad, o types.Order) error {
	if hasDatePassed(o.LaunchDate, launchpad.Location) {
		return types.NewErrInvalidData("launch date has passed")
	}
	firstDestination, err := s.launchpadFirstDestinationRepo.Get(ctx, o.LaunchpadID)
	if err != nil {
		return errors.Errorf(`no first destination for launchpad: id - %s`, o.LaunchpadID)
	}
	destinations, err := s.destinationRepo.ListSorted(ctx)
	if err != nil {
		return errors.Wrap(err, `failed to get destinations`)
	}
	if len(destinations) == 0 {
		return types.ErrFlightImpossible{}
	}
	destinationID, err := calculateDestinationForDate(o.LaunchDate, launchpad.Location, firstDestination, destinations)
	if err != nil {
		return err
	}
	if destinationID != o.DestinationID {
		return types.ErrFlightImpossible{}
	}
	return nil
}

func calculateDestinationForDate(requestedDate time.Time, location *time.Location, first types.LaunchpadFirstDestination, destinations []types.Destination) (string, error) {
	year, month, day := requestedDate.In(location).Date()
	requestedDateStartOfDay := time.Date(year, month, day, 0, 0, 0, 0, location)
	firstLaunchStartOfDay := time.Date(first.LocalYear, time.Month(first.LocalMonth), first.LocalDay, 0, 0, 0, 0, location)
	daysShift := int(requestedDateStartOfDay.Sub(firstLaunchStartOfDay).Hours() / 24)
	destinationsN := len(destinations)
	var firstDestinationOrder int
	var destinationFound bool
	for i, dest := range destinations {
		if dest.ID == first.DestinationID {
			destinationFound = true
			firstDestinationOrder = i + 1
			break
		}
	}
	if !destinationFound {
		return "", errors.Errorf(`first destination is not present in destinations: first - %+v, destinations: - %+v`, first, destinations)
	}
	destinationsShift := daysShift % destinationsN
	destinationOrder := firstDestinationOrder + destinationsShift
	if destinationOrder > destinationsN {
		destinationOrder = destinationOrder - destinationsN
	}
	return destinations[destinationOrder-1].ID, nil
}

func hasDatePassed(requestedLaunchDate time.Time, location *time.Location) bool {
	nowYear, nowMonth, nowDay := time.Now().In(location).Date()
	requestedYear, requestedMonth, requestedDay := requestedLaunchDate.In(location).Date()
	if requestedYear < nowYear {
		return true
	}
	if requestedYear > nowYear {
		return false
	}
	if requestedMonth < nowMonth {
		return true
	}
	if requestedMonth > nowMonth {
		return false
	}
	if requestedDay < nowDay {
		return true
	}
	return false
}

func (s *Orders) Get(ctx context.Context, id string) (types.Order, error) {
	return s.orderRepo.Get(ctx, id)
}

func (s *Orders) List(ctx context.Context, limit, offset int) ([]types.Order, error) {
	return s.orderRepo.List(ctx, limit, offset)
}

func (s *Orders) Delete(ctx context.Context, id string) error {
	return s.orderRepo.Delete(ctx, id)
}
