package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)

	router.Get("/", handler.HomeHandler)
	router.Post("/register", userHandler.Register)

	// room handlers
	router.Route("/rooms", func(router chi.Router) {
		router.Use(userMiddleware.AuthUserCtx)

		router.Post("/", roomHandler.Create)
		router.Get("/templates", roomHandler.VoteTemplates)

		router.Route("/{roomId}", func(router chi.Router) {
			router.Use(roomMiddleware.RoomCtx)

			router.Get("/", roomHandler.Show)

			router.Group(func(router chi.Router) {
				router.Use(roomMiddleware.RoomParticipant)

				router.Post("/leave", roomHandler.Leave)
				router.Post("/vote", roomHandler.Vote)
				router.Get("/subscribe", roomHandler.Subscribe)
			})

			router.Group(func(router chi.Router) {
				router.Use(roomMiddleware.RoomOwner)

				router.Patch("/", roomHandler.Update)
				router.Post("/end", roomHandler.EndVote)
				router.Post("/reset", roomHandler.Reset)
			})
		})
	})

	// example route
	router.Handle("/example/*", http.StripPrefix("/example/", http.FileServer(http.Dir("example"))))

	port := os.Getenv("POPLAN_PORT")
	if port == "" {
		port = "80"
	}
	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(err)
	}
}
