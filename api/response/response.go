package response

import "github.com/violarium/poplan/room"

type Home struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Message struct {
	Message string `json:"message"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Registration struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type Room struct {
	Id     string      `json:"id"`
	Name   string      `json:"name"`
	Status room.Status `json:"status"`
	Seats  []Seat      `json:"seats"`
}

type Seat struct {
	User  User `json:"user"`
	Vote  uint `json:"vote"`
	Voted bool `json:"voted"`
}
