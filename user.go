package main

import (
	"github.com/google/uuid"
)

type User struct {
	id   string
	name string
}

func NewUser(name string) *User {
	id := uuid.NewString()
	return &User{id: id, name: name}
}
