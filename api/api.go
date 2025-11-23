package api

import (
	"database/sql"

	"github.com/rafaeldepontes/go-full-crud/internal/repository"
	"github.com/rafaeldepontes/go-full-crud/internal/usecase"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	*usecase.UserHandler
}

func Init() (*sql.DB, error) {
	db, err := repository.Open()

	if err != nil {
		log.Error(err)
		return nil, err
	}
	return db, nil
}
