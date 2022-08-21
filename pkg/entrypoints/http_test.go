package entrypoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/leveldorado/space-trouble/pkg/tools/logger"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/leveldorado/space-trouble/pkg/types"
	"github.com/pkg/errors"
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
		BirthdayYear:  1990,
		BirthdayDay:   10,
		BirthdayMonth: 13,
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

func TestCreateOrderInvalid(t *testing.T) {
	o := types.Order{}
	b := &bytes.Buffer{}
	require.NoError(t, json.NewEncoder(b).Encode(o))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", b)
	resp := httptest.NewRecorder()
	NewHTTPEntry(nil, &logrus.Logger{}).GetHandler().ServeHTTP(resp, req)
	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateOrderFails(t *testing.T) {
	errorCases := []error{
		types.ErrFlightImpossible{},
		types.ErrInvalidData{},
		types.ErrDuplicatedOrder{},
		errors.New("fail"),
	}
	expectedCodes := []int{
		http.StatusNotAcceptable,
		http.StatusBadRequest,
		http.StatusConflict,
		http.StatusInternalServerError,
	}
	var orders []types.Order
	s := &mockOrdersService{}
	for _, err := range errorCases {
		order := types.Order{
			FirstName:     gofakeit.FirstName(),
			LastName:      gofakeit.LastName(),
			Gender:        gofakeit.Gender(),
			BirthdayYear:  1990,
			BirthdayDay:   10,
			BirthdayMonth: 13,
			LaunchpadID:   uuid.New().String(),
			DestinationID: uuid.New().String(),
			LaunchDate:    time.Now().UTC(),
		}
		s.On("Create", mock.Anything, order).Return("", err)
		orders = append(orders, order)
	}
	h := NewHTTPEntry(s, logger.New()).GetHandler()
	for i, order := range orders {
		b := &bytes.Buffer{}
		require.NoError(t, json.NewEncoder(b).Encode(order))
		req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", b)
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		require.Equal(t, expectedCodes[i], resp.Code)
	}
	s.AssertExpectations(t)
}

func TestOrdersList(t *testing.T) {
	orders := []types.Order{
		{
			ID:            uuid.New().String(),
			FirstName:     gofakeit.FirstName(),
			LastName:      gofakeit.LastName(),
			Gender:        gofakeit.Gender(),
			BirthdayYear:  1990,
			BirthdayDay:   10,
			BirthdayMonth: 13,
			LaunchpadID:   uuid.New().String(),
			DestinationID: uuid.New().String(),
			LaunchDate:    gofakeit.Date(),
		},
	}
	limit, offset := 10, 30
	s := &mockOrdersService{}
	s.On("ListSorted", mock.Anything, limit, offset).Return(orders, nil)

	h := NewHTTPEntry(s, &logrus.Logger{}).GetHandler()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/orders?limit=%d&offset=%d", limit, offset), nil)
	resp := httptest.NewRecorder()

	h.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var resultOrders []types.Order
	require.NoError(t, json.Unmarshal([]byte(gjson.GetBytes(resp.Body.Bytes(), "docs").Raw), &resultOrders))
	require.Equal(t, orders, resultOrders)
	require.Equal(t, limit, int(gjson.GetBytes(resp.Body.Bytes(), "limit").Int()))
	require.Equal(t, offset, int(gjson.GetBytes(resp.Body.Bytes(), "offset").Int()))

	s.AssertExpectations(t)
}

func TestGetOrder(t *testing.T) {
	order := types.Order{
		ID:            uuid.New().String(),
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		Gender:        gofakeit.Gender(),
		BirthdayYear:  1990,
		BirthdayDay:   10,
		BirthdayMonth: 13,
		LaunchpadID:   uuid.New().String(),
		DestinationID: uuid.New().String(),
		LaunchDate:    gofakeit.Date(),
	}
	s := &mockOrdersService{}
	s.On("Get", mock.Anything, order.ID).Return(order, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+order.ID, nil)
	resp := httptest.NewRecorder()

	NewHTTPEntry(s, &logrus.Logger{}).GetHandler().ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var resultOrder types.Order
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &resultOrder))
	require.Equal(t, order, resultOrder)

	s.AssertExpectations(t)
}

func TestDeleteOrder(t *testing.T) {
	id := uuid.New().String()
	s := &mockOrdersService{}
	s.On("Delete", mock.Anything, id).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/orders/"+id, nil)
	resp := httptest.NewRecorder()

	NewHTTPEntry(s, &logrus.Logger{}).GetHandler().ServeHTTP(resp, req)

	require.Equal(t, http.StatusNoContent, resp.Code)

	s.AssertExpectations(t)
}
