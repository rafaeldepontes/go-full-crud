package database

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

func Open() (*sql.DB, error) {
	conn, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}

	println("Database successfully connected...")
	return conn, nil
}
