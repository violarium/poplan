package user

import (
	"sync"
)

// todo: Register limit for users

type Registry struct {
	sync.RWMutex
	roster map[string]*User
}

func NewUserRegistry() *Registry {
	return &Registry{roster: make(map[string]*User)}
}

func (registry *Registry) Register(user *User) string {
	registry.Lock()
	defer registry.Unlock()
	registry.roster[user.Id] = user

	return user.Id
}

func (registry *Registry) Find(token string) (*User, bool) {
	registry.RLock()
	defer registry.RUnlock()
	user, ok := registry.roster[token]

	return user, ok
}
