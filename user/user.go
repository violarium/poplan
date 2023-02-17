package user

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

func (u *User) Id() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}
