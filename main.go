package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func setJsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type Home struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Register struct {
	Name string `json:"name"`
}

type UserNote struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Registration struct {
	UserNote UserNote `json:"user"`
	Token    string   `json:"token"`
}

func handleHome(w http.ResponseWriter, _ *http.Request) {
	home := Home{Title: "PoPlan", Description: "Planning Poker"}
	err := json.NewEncoder(w).Encode(home)
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	userRegistry := NewUserRegistry()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(setJsonContentType)

	router.Get("/", handleHome)

	router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		var register Register
		if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
			log.Println(err)
			return
		}

		if register.Name == "" {
			// todo: return error 400 or generate name
			return
		}

		user := NewUser(register.Name)
		token := userRegistry.register(user)

		err := json.NewEncoder(w).Encode(Registration{
			UserNote: UserNote{
				Id:   user.id,
				Name: user.name,
			},
			Token: token,
		})
		if err != nil {
			log.Println(err)
			return
		}
	})

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
