package handler

import (
	"net/http"

	"github.com/violarium/poplan/api"
	"github.com/violarium/poplan/api/response"
)

func HomeHandler(w http.ResponseWriter, _ *http.Request) {
	home := response.Home{Title: "PoPlan", Description: "Planning Poker"}
	api.SendResponse(w, home, http.StatusOK)
}
