package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/go-full-crud/api"
	"github.com/rafaeldepontes/go-full-crud/internal/handler"
	"github.com/rafaeldepontes/go-full-crud/internal/repository"
	"github.com/rafaeldepontes/go-full-crud/internal/usecase"
	log "github.com/sirupsen/logrus"
)

const banner = `
                		Go CRUD Service                 		  
 -------------------------------------------------------------------------
   Status  : Starting up
   Version : 1.1.0
   Go      : runtime.GoVersion()

   # if you're not using ".env" the listen should be "localhost:8000"
   Listen : Check your .env
 -------------------------------------------------------------------------
`

func main() {
	godotenv.Load(".env", ".env.example")
	log.SetReportCaller(true)

	db, err := api.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userHandler := usecase.NewUserHandler(userRepo)
	app := &api.Application{UserHandler: userHandler}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	recoveryDB := func(ctx context.Context) {
		_ = db.Close()

		newDB, err := api.Init()
		if err != nil {
			log.Errorf("[ERROR] Database recovery failed: %v", err)
			return
		}

		*db = *newDB
		log.Info("[INFO] DB recovered successfully")
	}

	go api.HealthCheck[*sql.DB](os.Getenv("DATABASE"), app.IsDbOk, ctx, 10*time.Second, recoveryDB)

	var r *chi.Mux = chi.NewRouter()
	handler.Handler(r, app)

	println(banner)

	var address string = "localhost:8000"
	if isTest := os.Getenv("ADDRESS"); isTest != "" {
		address = isTest
	}

	err = http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal(err)
	}
}
