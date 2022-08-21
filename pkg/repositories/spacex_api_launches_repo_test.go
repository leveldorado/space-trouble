package repositories

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

/*
Normally for testing external API in unit tests - should be used mock of http.Transport.

	so tests do not exceed requests quota
*/
func TestSpaceXAPILaunchesRepo_CheckLaunches(t *testing.T) {
	r := NewSpaceXAPILaunchesRepo(http.DefaultClient)
	launchpad := "5e9e4502f509092b78566f87"
	busyDay := time.Date(2022, 8, 26, 0, 0, 0, 0, time.UTC)
	notBusyDay := time.Date(2022, 8, 28, 0, 0, 0, 0, time.UTC)

	exists, err := r.CheckLaunches(context.TODO(), launchpad, busyDay)
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = r.CheckLaunches(context.TODO(), launchpad, notBusyDay)
	require.NoError(t, err)
	require.False(t, exists)
}
