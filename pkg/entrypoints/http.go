package entrypoints

import (
	"context"
	"encoding/json"
	"net/http"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/leveldorado/space-trouble/pkg/types"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type ordersService interface {
	Create(ctx context.Context, o types.Order) (string, error)
	Get(ctx context.Context, id string) (types.Order, error)
	List(ctx context.Context, limit, offset int64) ([]types.Order, error)
	Delete(ctx context.Context, id string) error
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
	r.Use(chiprometheus.NewMiddleware("space-trouble"))
	r.Get("/health", func(wr http.ResponseWriter, _ *http.Request) {
		wr.WriteHeader(http.StatusOK)
	})
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/orders", func(r chi.Router) {
			r.Post("/", e.createOrder)
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

func (e *HTTPEntry) respond(ctx context.Context, resp interface{}, err error, successCode int, wr http.ResponseWriter) {
	if err != nil {
		e.respondError(ctx, err, wr)
		return
	}
	if resp == nil {
		wr.WriteHeader(successCode)
		return
	}
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
