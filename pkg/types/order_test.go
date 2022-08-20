package types

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOrder_Validate(t *testing.T) {
	o := Order{}
	require.Error(t, o.Validate())
	o.FirstName = gofakeit.FirstName()
	require.Error(t, o.Validate())
	o.LastName = gofakeit.LastName()
	require.Error(t, o.Validate())
	o.Birthday = time.Date(2000, 10, 11, 0, 0, 0, 0, time.UTC)
	require.Error(t, o.Validate())
	o.LaunchpadID = uuid.New().String()
	require.Error(t, o.Validate())
	o.DestinationID = uuid.New().String()
	require.NoError(t, o.Validate())
}
