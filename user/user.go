package user

import (
	"github.com/google/uuid"
)

type User struct {
	Id   string
	Name string
}

func NewUser(name string) *User {
	id := uuid.NewString()
	return &User{Id: id, Name: name}
}
