package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/violarium/poplan/api"
	"github.com/violarium/poplan/api/request"
	"github.com/violarium/poplan/api/response"
	"github.com/violarium/poplan/user"
)

func GetRegisterHandler(userRegistry *user.Registry) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var register request.Register

		{
			err := json.NewDecoder(r.Body).Decode(&register)
			if err != nil || register.Name == "" {
				api.SendMessage(w, `"Name"" is required`, http.StatusUnprocessableEntity)
				return
			}
		}

		newUser := user.NewUser(register.Name)
		token := userRegistry.Register(newUser)
		registration := response.Registration{
			User: response.User{
				Id:   newUser.Id,
				Name: newUser.Name,
			},
			Token: token,
		}

		if err := json.NewEncoder(w).Encode(registration); err != nil {
			log.Println(err)
			return
		}
	}
}
