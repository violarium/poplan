package room

import (
	"errors"
	"sync"
)

type Registry struct {
	mu     sync.RWMutex
	roster map[string]*Room
}

func NewRegistry() *Registry {
	return &Registry{roster: make(map[string]*Room)}
}

func (registry *Registry) Add(room *Room) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if _, collision := registry.roster[room.Id()]; collision {
		return errors.New("room with such id already exists")
	}

	registry.roster[room.Id()] = room
	return nil
}

func (registry *Registry) Find(id string) (*Room, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	room, ok := registry.roster[id]

	return room, ok
}
