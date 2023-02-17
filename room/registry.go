package room

import (
	"sync"
)

type Registry struct {
	mu     sync.RWMutex
	roster map[string]*Room
}

func NewRegistry() *Registry {
	return &Registry{roster: make(map[string]*Room)}
}

func (registry *Registry) Add(room *Room) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.roster[room.Id()] = room
}

func (registry *Registry) Find(id string) (*Room, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	room, ok := registry.roster[id]

	return room, ok
}
