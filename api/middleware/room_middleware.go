package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/violarium/poplan/api"
	"github.com/violarium/poplan/room"
)

type RoomMiddleware struct {
	roomRegistry *room.Registry
}

func NewRoomMiddleware(roomRegistry *room.Registry) *RoomMiddleware {
	return &RoomMiddleware{roomRegistry: roomRegistry}
}

func (m *RoomMiddleware) RoomCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roomId := chi.URLParam(r, "roomId")
		foundRoom, ok := m.roomRegistry.Find(roomId)
		if !ok {
			api.SendMessage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "room", foundRoom)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *RoomMiddleware) RoomOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authUser, authUserOk := api.GetAuthUser(r)
		currentRoom, currentRoomOk := api.GetCurrentRoom(r)
		if !authUserOk || !currentRoomOk || authUser != currentRoom.Owner() {
			api.SendMessage(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
		next.ServeHTTP(w, r)
	})
}

func (m *RoomMiddleware) RoomParticipant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authUser, authUserOk := api.GetAuthUser(r)
		currentRoom, currentRoomOk := api.GetCurrentRoom(r)
		if !authUserOk || !currentRoomOk || !currentRoom.HasParticipant(authUser) {
			api.SendMessage(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
		next.ServeHTTP(w, r)
	})
}
