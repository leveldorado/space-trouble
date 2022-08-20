package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"golang.org/x/net/context"

	"github.com/leveldorado/space-trouble/pkg/entrypoints"
	"github.com/leveldorado/space-trouble/pkg/services"
	"github.com/leveldorado/space-trouble/pkg/tools/logger"
)

func main() {
	s := services.NewOrders()
	log := logger.New()

	h := entrypoints.NewHTTPEntry(s, log).GetHandler()

	httpS := &http.Server{
		Addr:         ":8000",
		Handler:      h,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("server started")
		if err := httpS.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithField("err", err.Error()).Error("failed to listen and serve http")
			signalChan <- syscall.SIGTERM
		}
	}()

	<-signalChan
	log.Info("server exiting")
	if err := httpS.Shutdown(context.Background()); err != nil {
		log.WithField("err", err.Error()).Error("failed to shutdown server")
		return
	}
	log.Info("BYE!")
}
