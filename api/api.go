package api

import (
	"context"
	"database/sql"
	"time"

	"github.com/rafaeldepontes/go-full-crud/internal/database"
	"github.com/rafaeldepontes/go-full-crud/internal/usecase"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	*usecase.Service
}

func Init() (*sql.DB, error) {
	db, err := database.Open()

	if err != nil {
		log.Error(err)
		return nil, err
	}
	return db, nil
}

func (app *Application) IsDbOk() bool {
	return app.Service.Repository.Ping()
}

func HealthCheck[T any](healthCheckName string, isOk func() bool, ctx context.Context, checkTime time.Duration, tryRecover func(ctx context.Context)) {
	var ticker *time.Ticker = time.NewTicker(checkTime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Infof("[HEALTH CHECK] %s: Ending session...", healthCheckName)
			return
		case <-ticker.C:
			ok := false

			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Errorf("[HEALTH CHECK] %s: panic in isOk(): %v", healthCheckName, r)
					}
				}()

				ok = isOk()
			}()

			if !ok {
				log.Warnf("[HEALTH CHECK] %s: Unhealthy state detected, running recovery", healthCheckName)
				tryRecover(ctx)
			} else {
				log.Infof("[HEALTH CHECK] %s: System is ok", healthCheckName)
			}
		}
	}
}
