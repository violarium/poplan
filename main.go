package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/violarium/poplan/api/handler"
	"github.com/violarium/poplan/api/middleware"
	"github.com/violarium/poplan/room"
	"github.com/violarium/poplan/user"
)

func main() {
	userRegistry := user.NewRegistry()
	roomRegistry := room.NewRegistry()

	userHandler := handler.NewUserHandler(userRegistry)
	userMiddleware := middleware.NewUserMiddleware(userRegistry)

	roomHandler := handler.NewRoomHandler(roomRegistry)
	roomMiddleware := middleware.NewRoomMiddleware(roomRegistry)

	router := chi.NewRouter()
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)

	router.Get("/", handler.HomeHandler)
	router.Post("/register", userHandler.Register)

	// room handlers
	router.Route("/room", func(router chi.Router) {
		router.Use(userMiddleware.AuthUserCtx)

		router.Post("/", roomHandler.Create)

		router.Route("/{roomId}", func(router chi.Router) {
			router.Use(roomMiddleware.RoomCtx)

			router.Get("/", func(w http.ResponseWriter, r *http.Request) {
				// todo: show room but only if user is owner or participant
			})

			router.Patch("/update", func(w http.ResponseWriter, r *http.Request) {
				// todo: update room
			})

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
