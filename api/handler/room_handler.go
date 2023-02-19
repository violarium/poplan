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
	h.roomRegistry.Add(newRoom)

	h.sendRoomResponse(w, newRoom, http.StatusCreated)
}

func (h *RoomHandler) Show(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	if !currentRoomOk {
		return
	}

	h.sendRoomResponse(w, currentRoom, http.StatusOK)
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

	h.sendRoomResponse(w, currentRoom, http.StatusOK)
}

func (h *RoomHandler) Join(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	authUser, authUserOk := api.GetAuthUser(r)
	if !currentRoomOk || !authUserOk {
		return
	}

	currentRoom.Join(authUser)

	h.sendRoomResponse(w, currentRoom, http.StatusOK)
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

	h.sendRoomResponse(w, currentRoom, http.StatusOK)
}

func (h *RoomHandler) Reset(w http.ResponseWriter, r *http.Request) {
	currentRoom, currentRoomOk := api.GetCurrentRoom(r)
	if !currentRoomOk {
		return
	}
	currentRoom.Reset()

	h.sendRoomResponse(w, currentRoom, http.StatusOK)
}

func (h *RoomHandler) sendRoomResponse(w http.ResponseWriter, r *room.Room, httpStatus int) {
	seats := r.Seats()
	seatsResponse := make([]response.Seat, 0, len(seats))
	for _, s := range seats {
		seatsResponse = append(seatsResponse, response.Seat{
			User: response.User{
				Id:   s.User().Id(),
				Name: s.User().Name(),
			},
			Vote: response.Vote{
				Value: s.SecretVote().Value(),
				Type:  s.SecretVote().Type(),
			},
			Voted: s.Voted(),
			Owner: s.User() == r.Owner(),
		})
	}

	votes := r.VoteTemplate().Votes
	voteResponses := make([]response.Vote, 0, len(votes))
	for _, v := range votes {
		voteResponses = append(voteResponses, response.Vote{
			Value: v.Value(),
			Type:  v.Type(),
		})
	}

	roomResponse := response.Room{
		Id:     r.Id(),
		Name:   r.Name(),
		Status: r.Status(),
		Seats:  seatsResponse,
		VoteTemplate: response.VoteTemplate{
			Title: r.VoteTemplate().Title,
			Votes: voteResponses,
		},
	}

	api.SendResponse(w, roomResponse, httpStatus)
}
