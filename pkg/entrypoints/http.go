package entrypoints

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/leveldorado/space-trouble/pkg/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ordersService interface {
	Create(ctx context.Context, o types.Order) (string, error)
	Get(ctx context.Context, id string) (types.Order, error)
	List(ctx context.Context, limit, offset int) ([]types.Order, error)
	Delete(ctx context.Context, id string) error
	Destinations(ctx context.Context) ([]types.Destination, error)
}

type HTTPEntry struct {
	os  ordersService
	log logrus.FieldLogger
}

func NewHTTPEntry(os ordersService, log logrus.FieldLogger) *HTTPEntry {
	return &HTTPEntry{os: os, log: log}
}

func (e *HTTPEntry) GetHandler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Get("/health", func(wr http.ResponseWriter, _ *http.Request) {
		wr.WriteHeader(http.StatusOK)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/orders", func(r chi.Router) {
			r.Post("/", e.createOrder)
			r.Get("/", e.list)
			r.Get("/{id}", e.getOrder)
			r.Delete("/{id}", e.deleteOrder)
		})
		r.Route("/destinations", func(r chi.Router) {
			r.Get("/", e.destinations)
		})
	})
	return r
}

type createOrderResponse struct {
	ID string `json:"id"`
}

func (e *HTTPEntry) createOrder(wr http.ResponseWriter, r *http.Request) {
	o := types.Order{}
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		e.respondError(r.Context(), types.NewErrInvalidData(err.Error()), wr)
		return
	}
	if err := o.Validate(); err != nil {
		e.respondError(r.Context(), types.NewErrInvalidData(err.Error()), wr)
		return
	}
	id, err := e.os.Create(r.Context(), o)
	e.respond(r.Context(), createOrderResponse{ID: id}, err, http.StatusCreated, wr)
}

type paginationResult struct {
	Docs   interface{} `json:"docs"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

const (
	defaultOrdersLimit = 10
	maxOrdersLimit     = 1000
)

func (e *HTTPEntry) list(wr http.ResponseWriter, req *http.Request) {
	limit, offset, err := parseLimitOffset(req.URL.Query(), defaultOrdersLimit, maxOrdersLimit)
	if err != nil {
		e.respondError(req.Context(), err, wr)
		return
	}
	orders, err := e.os.List(req.Context(), limit, offset)
	e.respond(req.Context(), paginationResult{
		Docs:   orders,
		Limit:  limit,
		Offset: offset,
	}, err, http.StatusOK, wr)
}

func parseLimitOffset(values url.Values, defaultLimit, maxLimit int) (int, int, error) {
	limit, err := parseIntQueryParam(values, "limit", defaultOrdersLimit)
	if err != nil {
		return 0, 0, err
	}
	if limit > maxOrdersLimit {
		return 0, 0, types.NewErrInvalidData("limit param exceed max value " + strconv.Itoa(maxOrdersLimit))
	}
	offset, err := parseIntQueryParam(values, "offset", 0)
	if err != nil {
		return 0, 0, err
	}
	return limit, offset, err
}

func parseIntQueryParam(values url.Values, key string, defaultValue int) (int, error) {
	str := values.Get(key)
	if str == "" {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(str)
	if err != nil {
		return 0, types.NewErrInvalidData("invalid value " + str + " for key " + key)
	}
	return value, nil
}

func (e *HTTPEntry) getOrder(wr http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	order, err := e.os.Get(req.Context(), id)
	e.respond(req.Context(), order, err, http.StatusOK, wr)
}

func (e *HTTPEntry) deleteOrder(wr http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	err := e.os.Delete(req.Context(), id)
	e.respond(req.Context(), nil, err, http.StatusNoContent, wr)
}

func (e *HTTPEntry) respond(ctx context.Context, resp interface{}, err error, successCode int, wr http.ResponseWriter) {
	if err != nil {
		e.respondError(ctx, err, wr)
		return
	}
	if resp == nil {
		wr.WriteHeader(successCode)
		return
	}
	wr.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		e.respondError(ctx, errors.Wrapf(err, `failed to marshal response`), wr)
		return
	}
	wr.WriteHeader(successCode)
	_, err = wr.Write(data)
	if err != nil {
		e.log.WithField("err", err.Error()).
			WithField("data", string(data)).WithContext(ctx).
			Warn("failed to write response")
	}
}

type errorResponse struct {
	Message string `json:"message"`
}

func (e *HTTPEntry) respondError(ctx context.Context, err error, wr http.ResponseWriter) {
	logEntry := e.log.WithField("err", err.Error()).WithContext(ctx)
	resp := errorResponse{}
	code := http.StatusInternalServerError
	switch cause := errors.Cause(err).(type) {
	case types.ErrFlightImpossible:
		resp.Message = cause.Error()
		code = http.StatusNotAcceptable
	case types.ErrInvalidData:
		resp.Message = cause.Error()
		code = http.StatusBadRequest
	case types.ErrDuplicatedOrder:
		resp.Message = cause.Error()
		code = http.StatusConflict
	default:
		resp.Message = err.Error()
	}
	logEntry = logEntry.WithField("code", code)
	if code != http.StatusInternalServerError {
		logEntry.Warn("non success response")
	} else {
		logEntry.Error("failed to handle request")
	}
	wr.WriteHeader(code)
	wr.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		e.log.WithField("err", err.Error()).
			WithField("response", resp).WithContext(ctx).
			Error("failed to marshal error response")
		return
	}
	_, err = wr.Write(data)
	if err != nil {
		e.log.WithField("err", err.Error()).
			WithField("data", string(data)).WithContext(ctx).
			Warn("failed to write error response")
		return
	}
}

func (e *HTTPEntry) destinations(wr http.ResponseWriter, req *http.Request) {
	destinations, err := e.os.Destinations(req.Context())
	e.respond(req.Context(), destinations, err, http.StatusOK, wr)
}
