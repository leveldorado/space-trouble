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
	}, pad)
}
