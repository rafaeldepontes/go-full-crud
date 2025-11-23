package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/go-full-crud/api"
	"github.com/rafaeldepontes/go-full-crud/internal/handler"
	"github.com/rafaeldepontes/go-full-crud/internal/repository"
	"github.com/rafaeldepontes/go-full-crud/internal/usecase"
	log "github.com/sirupsen/logrus"
)

const banner =
`
                		Go CRUD Service                 		  
 -------------------------------------------------------------------------
   Status  : Starting up
   Version : 1.0.0
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

	applicationContext := &api.Application{
		UserHandler: userHandler,
	}

	var r *chi.Mux = chi.NewRouter()
	handler.Handler(r, applicationContext)

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
