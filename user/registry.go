package user

import (
	"sync"
)

type Registry struct {
	mu     sync.RWMutex
	roster map[string]*User
}

func NewRegistry() *Registry {
	return &Registry{roster: make(map[string]*User)}
}

func (registry *Registry) Register(user *User) string {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.roster[user.Id()] = user

	return user.Id()
}

func (registry *Registry) Find(token string) (*User, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	user, ok := registry.roster[token]

	return user, ok
}
