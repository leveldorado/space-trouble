package types

import (
	"testing"

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
	o.BirthdayYear = 2000
	o.BirthdayMonth = 10
	o.BirthdayDay = 23
	require.Error(t, o.Validate())
	o.LaunchpadID = uuid.New().String()
	require.Error(t, o.Validate())
	o.DestinationID = uuid.New().String()
	require.NoError(t, o.Validate())
}
