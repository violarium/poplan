package user

import (
	"errors"
	"sync"
)

type Registry struct {
	mu     sync.RWMutex
	roster map[string]*User
}

func NewRegistry() *Registry {
	return &Registry{roster: make(map[string]*User)}
}

func (registry *Registry) Register(user *User) (string, error) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if _, collision := registry.roster[user.Id()]; collision {
		return "", errors.New("user with such id already exists")
	}
	registry.roster[user.Id()] = user

	return user.Id(), nil
}

func (registry *Registry) Find(token string) (*User, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	user, ok := registry.roster[token]

	return user, ok
}
