package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/violarium/poplan/api/response"
)

func HomeHandler(w http.ResponseWriter, _ *http.Request) {
	home := response.Home{Title: "PoPlan", Description: "Planning Poker"}
	err := json.NewEncoder(w).Encode(home)
	if err != nil {
		log.Println(err)
		return
	}
}
