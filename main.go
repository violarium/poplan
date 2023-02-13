package main

import (
	"encoding/json"
	"github.com/violarium/poplan/api"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var userRegistry = NewUserRegistry()

func main() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(api.SetJsonContentType)

	router.Get("/", handleHome)
	router.Post("/register", handleRegister)

	// todo: use context to pass room and authorized user

	router.Post("/room", func(w http.ResponseWriter, r *http.Request) {
		// todo: create room and register to room list
	})

	router.Route("/room/{id}", func(router chi.Router) {
		router.Put("/vote", func(w http.ResponseWriter, r *http.Request) {
			// todo: user votes
		})

		router.Post("/reveal", func(w http.ResponseWriter, r *http.Request) {
			// todo: creator reveals
		})

		router.Post("/reset", func(w http.ResponseWriter, r *http.Request) {
			// todo: creator resets room
		})
	}) // todo: register middleware, user has to be authorized

	// todo: use env to set port
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func handleHome(w http.ResponseWriter, _ *http.Request) {
	home := api.Home{Title: "PoPlan", Description: "Planning Poker"}
	err := json.NewEncoder(w).Encode(home)
	if err != nil {
		log.Println(err)
		return
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var register api.Register

	{
		err := json.NewDecoder(r.Body).Decode(&register)
		if err != nil || register.Name == "" {
			api.SendMessage(w, `"name"" is required`, http.StatusUnprocessableEntity)
			return
		}
	}

	user := NewUser(register.Name)
	token := userRegistry.register(user)
	registration := api.Registration{
		User: api.User{
			Id:   user.id,
			Name: user.name,
		},
		Token: token,
	}

	if err := json.NewEncoder(w).Encode(registration); err != nil {
		log.Println(err)
		return
	}
}
