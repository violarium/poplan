package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/violarium/poplan/api"
	"github.com/violarium/poplan/api/request"
	"github.com/violarium/poplan/api/response"
	"github.com/violarium/poplan/room"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
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
		if err != nil {
			api.SendMessage(w, `Json body required`, http.StatusUnprocessableEntity)
		}
		if createRoom.Name == "" {
			api.SendMessage(w, `"name" is required`, http.StatusUnprocessableEntity)
			return
		}
	}

	if createRoom.VoteTemplate >= len(room.DefaultVoteTemplates) {
		api.SendMessage(w, `Select valid "voteTemplate"`, http.StatusUnprocessableEntity)
		return
	}
	voteTemplate := room.DefaultVoteTemplates[createRoom.VoteTemplate]

	newRoom := room.NewRoom(authUser, createRoom.Name, voteTemplate)
	if err := h.roomRegistry.Add(newRoom); err != nil {
		api.SendMessage(w, "Can't create room, try again", http.StatusUnprocessableEntity)
		return
	}

	api.SendResponse(w, response.NewRoom(newRoom), http.StatusCreated)
}

func (h *RoomHandler) Show(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	if !currentRoomOk {
		return
	}

	api.SendResponse(w, response.NewRoom(currentRoom), http.StatusOK)
}

func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	if !currentRoomOk {
		return
	}

	var updateRoom request.UpdateRoom
	{
		err := json.NewDecoder(r.Body).Decode(&updateRoom)
		if err == nil && updateRoom.Name != "" {
			currentRoom.SetName(updateRoom.Name)
		}
	}

	api.SendResponse(w, response.NewRoom(currentRoom), http.StatusOK)
}

func (h *RoomHandler) Join(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	authUser, authUserOk := api.GetAuthUser(r)
	if !currentRoomOk || !authUserOk {
		return
	}

	currentRoom.Join(authUser)

	api.SendResponse(w, response.NewRoom(currentRoom), http.StatusOK)
}

func (h *RoomHandler) Leave(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	authUser, authUserOk := api.GetAuthUser(r)
	if !currentRoomOk || !authUserOk {
		return
	}

	currentRoom.Leave(authUser)

	api.SendMessage(w, "Room left", http.StatusOK)
}

func (h *RoomHandler) Vote(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	authUser, authUserOk := api.GetAuthUser(r)
	if !currentRoomOk || !authUserOk {
		return
	}

	var voteRequest request.Vote
	{
		err := json.NewDecoder(r.Body).Decode(&voteRequest)
		if err != nil {
			api.SendMessage(w, `"value" is required`, http.StatusUnprocessableEntity)
			return
		}
	}

	currentRoom.Vote(authUser, voteRequest.Value)

	api.SendMessage(w, "Voted", http.StatusOK)
}

func (h *RoomHandler) EndVote(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	if !currentRoomOk {
		return
	}
	currentRoom.EndVote()

	api.SendResponse(w, response.NewRoom(currentRoom), http.StatusOK)
}

func (h *RoomHandler) Reset(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	if !currentRoomOk {
		return
	}
	currentRoom.Reset()

	api.SendResponse(w, response.NewRoom(currentRoom), http.StatusOK)
}

func (h *RoomHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	authUser, authUserOk := api.GetAuthUser(r)
	if !currentRoomOk || !authUserOk {
		return
	}

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		api.SendMessage(w, "Can't establish websocket", http.StatusInternalServerError)
		return
	}
	defer (func() {
		if closeErr := c.Close(websocket.StatusInternalError, ""); closeErr != nil {
			log.Println("Close error", closeErr)
		}
	})()
	ctx := c.CloseRead(r.Context())

	subscriber, subscriberError := currentRoom.Subscribe(authUser)
	if subscriberError != nil {
		log.Println(err)
		api.SendMessage(w, "Unable to subscribe", http.StatusInternalServerError)
		return
	}
	defer currentRoom.Unsubscribe(subscriber)

	pingTicker := time.NewTicker(time.Second * 5)
	defer pingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			if closeErr := c.Close(websocket.StatusNormalClosure, ""); closeErr != nil {
				log.Println("Close error", closeErr)
			}
			return
		case <-pingTicker.C:
			if pingErr := c.Ping(ctx); pingErr != nil {
				log.Println("Ping error", pingErr)
				return
			}
		case _, ok := <-subscriber.Notifications:
			if !ok {
				log.Println("Room closed a channel")
				return
			}
			message := response.NewRoom(currentRoom)
			if writeErr := wsjson.Write(ctx, c, message); writeErr != nil {
				log.Println("Write error", writeErr)
				return
			}
		}
	}
}
