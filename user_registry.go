package main

import "sync"

// todo: register limit for users

type UserRegistry struct {
	sync.RWMutex
	roster map[string]*User
}

func NewUserRegistry() *UserRegistry {
	return &UserRegistry{roster: make(map[string]*User)}
}

func (registry *UserRegistry) register(user *User) string {
	registry.Lock()
	defer registry.Unlock()
	registry.roster[user.id] = user

	return user.id
}

func (registry *UserRegistry) find(token string) (*User, bool) {
	registry.RLock()
	defer registry.RUnlock()
	user, ok := registry.roster[token]

	return user, ok
}
