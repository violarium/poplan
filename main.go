package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/violarium/poplan/api"
	"github.com/violarium/poplan/user"
)

var userRegistry = user.NewUserRegistry()

func main() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(api.SetContentType("application/json"))

	router.Get("/", handleHome)
	router.Post("/register", handleRegister)

	// todo: use context to pass room and authorized user

	router.Post("/room", func(w http.ResponseWriter, r *http.Request) {
		// todo: create room and Register to room list
	})

	router.Route("/room/{Id}", func(router chi.Router) {
		router.Put("/vote", func(w http.ResponseWriter, r *http.Request) {
			// todo: user votes
		})

		router.Post("/reveal", func(w http.ResponseWriter, r *http.Request) {
			// todo: creator reveals
		})

		router.Post("/reset", func(w http.ResponseWriter, r *http.Request) {
			// todo: creator resets room
		})
	}) // todo: Register middleware, user has to be authorized

	port := os.Getenv("POPLAN_PORT")
	if port == "" {
		port = "80"
	}
	if err := http.ListenAndServe(":"+port, router); err != nil {
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
			api.SendMessage(w, `"Name"" is required`, http.StatusUnprocessableEntity)
			return
		}
	}

	u := user.NewUser(register.Name)
	token := userRegistry.Register(u)
	registration := api.Registration{
		User: api.User{
			Id:   u.Id,
			Name: u.Name,
		},
		Token: token,
	}

	if err := json.NewEncoder(w).Encode(registration); err != nil {
		log.Println(err)
		return
	}
}
