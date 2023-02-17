package handler

import (
	"encoding/json"
	"net/http"

	"github.com/violarium/poplan/api"
	"github.com/violarium/poplan/api/request"
	"github.com/violarium/poplan/api/response"
	"github.com/violarium/poplan/room"
)

type RoomHandler struct {
	roomRegistry *room.Registry
}

func NewRoomHandler(roomRegistry *room.Registry) *RoomHandler {
	return &RoomHandler{roomRegistry: roomRegistry}
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	authUser, authUserOk := api.GetAuthUser(r)
	if !authUserOk {
		return
	}

	var createRoom request.CreateRoom
	{
		err := json.NewDecoder(r.Body).Decode(&createRoom)
		if err != nil || createRoom.Name == "" {
			api.SendMessage(w, `"Name" is required`, http.StatusUnprocessableEntity)
			return
		}
	}

	newRoom := room.NewRoom(authUser, createRoom.Name)
	h.roomRegistry.Add(newRoom)

	newRoomResponse := response.Room{
		Id:   newRoom.Id(),
		Name: newRoom.Name(),
	}
	api.SendResponse(w, newRoomResponse, http.StatusCreated)
}
