package entrypoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/leveldorado/space-trouble/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestCreateOrderSuccess(t *testing.T) {
	o := types.Order{
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		Gender:        gofakeit.Gender(),
		Birthday:      time.Now().UTC(),
		LaunchpadID:   uuid.New().String(),
		DestinationID: uuid.New().String(),
		LaunchDate:    time.Now().UTC(),
	}
	id := uuid.New().String()
	os := &mockOrdersService{}
	os.On("Create", mock.Anything, o).Return(id, nil)

	h := NewHTTPEntry(os, &logrus.Logger{}).GetHandler()

	b := &bytes.Buffer{}
	require.NoError(t, json.NewEncoder(b).Encode(o))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", b)
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)
	require.Equal(t, id, gjson.GetBytes(resp.Body.Bytes(), "id").String())

	os.AssertExpectations(t)
}
