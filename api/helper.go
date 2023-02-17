package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/violarium/poplan/api/response"
	"github.com/violarium/poplan/room"
	"github.com/violarium/poplan/user"
)

func SendResponse(w http.ResponseWriter, response any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
	}
}

func SendMessage(w http.ResponseWriter, text string, status int) {
	msg := response.Message{Message: text}
	SendResponse(w, msg, status)
}

func GetAuthUser(r *http.Request) (*user.User, bool) {
	u, ok := r.Context().Value("authUser").(*user.User)

	return u, ok
}

func GetRoom(r *http.Request) (*room.Room, bool) {
	foundRoom, ok := r.Context().Value("room").(*room.Room)

	return foundRoom, ok
}
