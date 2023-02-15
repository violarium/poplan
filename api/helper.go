package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/violarium/poplan/api/response"
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
