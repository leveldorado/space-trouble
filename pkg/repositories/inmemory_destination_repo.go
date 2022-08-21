package repositories

import (
	"context"

	"github.com/leveldorado/space-trouble/pkg/types"
)

/*
InMemoryDestinationsRepo

		implementation for destinations repo.
	    for simplifying purposes without storing
*/
type InMemoryDestinationsRepo struct {
	destinations []types.Destination
}

func NewInMemoryDestinationsRepo() *InMemoryDestinationsRepo {
	return &InMemoryDestinationsRepo{
		destinations: []types.Destination{
			{
				ID:   "1",
				Name: "Mars",
			},
			{
				ID:   "2",
				Name: "IO",
			},
			{
				ID:   "3",
				Name: "Venus",
			},
			{
				ID:   "4",
				Name: "Jupiter",
			},
			{
				ID:   "5",
				Name: "Moon",
			},
			{
				ID:   "6",
				Name: "Neptune",
			},
			{
				ID:   "7",
				Name: "Pluto",
			},
		},
	}
}

func (r *InMemoryDestinationsRepo) ListSorted(_ context.Context) ([]types.Destination, error) {
	return r.destinations, nil
}
