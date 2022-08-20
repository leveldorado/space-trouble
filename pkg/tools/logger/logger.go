package logger

import (
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func New() *logrus.Logger {
	logger := logrus.New()
	logger.AddHook(requestIDHook{})
	return logger
}

type requestIDHook struct{}

func (requestIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (requestIDHook) Fire(e *logrus.Entry) error {
	requestID := middleware.GetReqID(e.Context)
	if requestID == "" {
		return nil
	}
	if e.Data == nil {
		e.Data = logrus.Fields{}
	}
	e.Data["request_id"] = requestID
	return nil
}
