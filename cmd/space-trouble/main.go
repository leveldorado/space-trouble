package main

import (
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leveldorado/space-trouble/pkg/repositories"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"golang.org/x/net/context"

	"github.com/leveldorado/space-trouble/pkg/entrypoints"
	"github.com/leveldorado/space-trouble/pkg/services"
	"github.com/leveldorado/space-trouble/pkg/tools/logger"
)

func main() {
	log := logger.New()
	conn := mustGetPostgresDB(log)
	cl := &http.Client{
		Timeout: time.Second,
	}

	s := services.NewOrders(
		repositories.NewPostgreSQLOrdersRepo(conn, log),
		repositories.NewSpaceXAPILaunchpadsRepo(cl),
		repositories.NewInMemoryDestinationsRepo(),
		repositories.NewInMemoryLaunchpadFirstDestinationRepo(),
		repositories.NewSpaceXAPILaunchesRepo(cl),
	)

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

func mustGetPostgresDB(log logrus.FieldLogger) *sql.DB {
	url := os.Getenv("POSTGRESQl_URL")
	db, err := repositories.GetPostgresqlConn(url)
	if err != nil {
		log.WithField("err", err.Error()).Fatal("failed to obtain postgres conn")
	}
	return db
}
