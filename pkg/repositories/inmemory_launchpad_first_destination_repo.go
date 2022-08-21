package repositories

import (
	"context"
	"fmt"
	"sync"

	"github.com/leveldorado/space-trouble/pkg/types"
)

/*
InMemoryLaunchpadFirstDestinationRepo

	implements LaunchpadFirstDestinationRepo
	does not interact with db for simplifying purpose
*/
type InMemoryLaunchpadFirstDestinationRepo struct {
	launchpads map[string]types.LaunchpadFirstDestination
	sync.RWMutex
}

func NewInMemoryLaunchpadFirstDestinationRepo() *InMemoryLaunchpadFirstDestinationRepo {
	return &InMemoryLaunchpadFirstDestinationRepo{launchpads: map[string]types.LaunchpadFirstDestination{}}
}

func (r *InMemoryLaunchpadFirstDestinationRepo) Set(doc types.LaunchpadFirstDestination) {
	r.Lock()
	r.launchpads[doc.LaunchpadID] = doc
	r.Unlock()
}

func (r *InMemoryLaunchpadFirstDestinationRepo) Get(_ context.Context, launchpad string) (types.LaunchpadFirstDestination, error) {
	r.RLock()
	defer r.RUnlock()
	doc, ok := r.launchpads[launchpad]
	fmt.Println(r.launchpads)
	if !ok {
		return types.LaunchpadFirstDestination{}, types.ErrNotFound{}
	}
	return doc, nil
}
