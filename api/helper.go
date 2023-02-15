package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/violarium/poplan/api/response"
)

func SendMessage(w http.ResponseWriter, text string, status int) {
	w.WriteHeader(status)
	msg := response.Message{Message: text}
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		log.Println(err)
	}
}
