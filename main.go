package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/violarium/poplan/api"
	"github.com/violarium/poplan/api/handler"
	"github.com/violarium/poplan/user"
)

func main() {
	userRegistry := user.NewRegistry()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(api.SetContentType("application/json"))

	router.Get("/", handler.HomeHandler)
	router.Post("/register", handler.GetRegisterHandler(userRegistry))

	router.Route("/room", func(router chi.Router) {
		router.Use(api.AuthUserCtx(userRegistry))
		router.Use(api.RequireAuthUser)

		router.Post("/", func(w http.ResponseWriter, r *http.Request) {
			if authUser, ok := r.Context().Value("authUser").(*user.User); !ok {
				// todo: create room
				fmt.Println(authUser)
			}
		})

		router.Route("/{id}", func(router chi.Router) {
			// todo: add middleware to get room by id
			// todo: add middleware to require room

			router.Post("/join", func(w http.ResponseWriter, r *http.Request) {
				// todo: user joins
			})

			router.Put("/vote", func(w http.ResponseWriter, r *http.Request) {
				// todo: user votes
			})

			// todo: add middlewares for ownership

			router.Post("/reveal", func(w http.ResponseWriter, r *http.Request) {
				// todo: creator reveals
			})

			router.Post("/reset", func(w http.ResponseWriter, r *http.Request) {
				// todo: creator resets room
			})
		})
	})

	port := os.Getenv("POPLAN_PORT")
	if port == "" {
		port = "80"
	}
	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(err)
	}
}
