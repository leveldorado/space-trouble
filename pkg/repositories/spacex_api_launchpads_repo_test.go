package repositories

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/leveldorado/space-trouble/pkg/types"
	"github.com/stretchr/testify/require"
)

/*
Normally for testing external API in unit tests - should be used mock of http.Transport.

	so tests do not exceed requests quota
*/
func TestSpaceXAPILaunchpadsRepo_Get(t *testing.T) {
	r := NewSpaceXAPILaunchpadsRepo(http.DefaultClient)
	padID := "5e9e4501f5090910d4566f83"
	losAngelesTimezone, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)

	pad, err := r.Get(context.TODO(), padID)

	require.NoError(t, err)
	require.Equal(t, types.Launchpad{
		ID:       "5e9e4501f5090910d4566f83",
		FullName: "Vandenberg Space Force Base Space Launch Complex 3W",
		Location: losAngelesTimezone,
		Status:   "retired",
	}, pad)
}

func TestSpaceXAPILaunchpadsRepo_List(t *testing.T) {
	newYorkTime, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	losAndgelestime, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)
	expected := []types.Launchpad{
		{
			ID:       "5e9e4501f509094ba4566f84",
			FullName: "Cape Canaveral Space Force Station Space Launch Complex 40",
			Location: newYorkTime,
		},
		{
			ID:       "5e9e4502f509094188566f88",
			FullName: "Kennedy Space Center Historic Launch Complex 39A",
			Location: newYorkTime,
		},
		{
			ID:       "5e9e4502f509092b78566f87",
			FullName: "Vandenberg Space Force Base Space Launch Complex 4E",
			Location: losAndgelestime,
		},
	}
	r := NewSpaceXAPILaunchpadsRepo(http.DefaultClient)
	resp, err := r.List(context.TODO())
	require.NoError(t, err)
	require.Equal(t, expected, resp)
}
