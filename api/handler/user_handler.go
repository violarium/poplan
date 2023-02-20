package handler

import (
	"encoding/json"
	"net/http"

	"github.com/violarium/poplan/api"
	"github.com/violarium/poplan/api/request"
	"github.com/violarium/poplan/api/response"
	"github.com/violarium/poplan/user"
)

type UserHandler struct {
	userRegistry *user.Registry
}

func NewUserHandler(userRegistry *user.Registry) *UserHandler {
	return &UserHandler{userRegistry: userRegistry}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var register request.Register

	{
		err := json.NewDecoder(r.Body).Decode(&register)
		if err != nil || register.Name == "" {
			api.SendMessage(w, `"name" is required`, http.StatusUnprocessableEntity)
			return
		}
	}

	newUser := user.NewUser(register.Name)
	token, err := h.userRegistry.Register(newUser)
	if err != nil {
		api.SendMessage(w, "Unable to register, try later", http.StatusUnprocessableEntity)
		return
	}

	registration := response.Registration{
		User: response.User{
			Id:   newUser.Id(),
			Name: newUser.Name(),
		},
		Token: token,
	}
	api.SendResponse(w, registration, http.StatusCreated)
}
