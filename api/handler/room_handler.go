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

	api.SendResponse(w, response.NewRoom(newRoom, authUser), http.StatusCreated)
}

func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	authUser, authUserOk := api.GetAuthUser(r)
	if !currentRoomOk || !authUserOk {
		return
	}

	var updateRoom request.UpdateRoom
	{
		err := json.NewDecoder(r.Body).Decode(&updateRoom)
		if err == nil && updateRoom.Name != "" {
			currentRoom.SetName(updateRoom.Name)
		}
	}

	api.SendResponse(w, response.NewRoom(currentRoom, authUser), http.StatusOK)
}

func (h *RoomHandler) Show(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	authUser, authUserOk := api.GetAuthUser(r)
	if !currentRoomOk || !authUserOk {
		return
	}

	currentRoom.Join(authUser)

	api.SendResponse(w, response.NewRoom(currentRoom, authUser), http.StatusOK)
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
	authUser, authUserOk := api.GetAuthUser(r)
	if !currentRoomOk || !authUserOk {
		return
	}
	currentRoom.EndVote()

	api.SendResponse(w, response.NewRoom(currentRoom, authUser), http.StatusOK)
}

func (h *RoomHandler) Reset(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	authUser, authUserOk := api.GetAuthUser(r)
	if !currentRoomOk || !authUserOk {
		return
	}
	currentRoom.Reset()

	api.SendResponse(w, response.NewRoom(currentRoom, authUser), http.StatusOK)
}

func (h *RoomHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	authUser, authUserOk := api.GetAuthUser(r)
	if !currentRoomOk || !authUserOk {
		return
	}

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
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

	// Async handle of notifications
	subscriber := room.NewSubscriber()

	// Middleware channel: subscriber notifications -> room changed
	roomChanged := make(chan bool, 1)

	go func() {
		for {
			// receive notification from room
			_, alive := <-subscriber.Notifications()
			if !alive {
				log.Println("No more notifications")
				close(roomChanged)
				return
			}
			log.Println("Notification received")

			// send event about room change
			select {
			case roomChanged <- true:
			default:
				// do nothing, already in process
				log.Println("Long connection, skip")
			}
		}
	}()

	// Subscribe
	currentRoom.Subscribe(authUser, subscriber)
	defer currentRoom.Unsubscribe(subscriber)
	log.Println("Subscribed")

	// Ticker to ping/pong websocket
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
		case _, alive := <-roomChanged:
			if !alive {
				log.Println("No more room changes")
				return
			}
			message := response.NewRoom(currentRoom, authUser)
			if writeErr := wsjson.Write(ctx, c, message); writeErr != nil {
				log.Println("Write error", writeErr)
				return
			}
		}
	}
}

func (h *RoomHandler) VoteTemplates(w http.ResponseWriter, _ *http.Request) {
	templateResponses := make([]response.VoteTemplate, 0, len(room.DefaultVoteTemplates))
	for _, template := range room.DefaultVoteTemplates {
		templateResponses = append(templateResponses, response.NewVoteTemplate(template))
	}

	api.SendResponse(w, response.VoteTemplateList{Templates: templateResponses}, http.StatusOK)
}
