package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func SendMessage(w http.ResponseWriter, text string, status int) {
	w.WriteHeader(status)
	msg := Message{Message: text}
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		log.Println(err)
	}
}
