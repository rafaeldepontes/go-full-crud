package handler

import (
	chimiddle "github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/rafaeldepontes/go-full-crud/api"
)

func Handler(r *chi.Mux, application *api.Application) {
	r.Use(chimiddle.StripSlashes)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/users", application.UserHandler.FindByUsername)
		r.Get("/users/{id}", application.UserHandler.FindUserById)

		r.Post("/users", application.UserHandler.Register)

		r.Put("/users/{id}", application.UserHandler.UpdateUserInfo)

		r.Delete("/users/{id}", application.UserHandler.DeleteUserById)
	})
}
