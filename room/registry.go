package room

import (
	"sync"
)

// todo: Register limit for rooms

type Registry struct {
	sync.RWMutex
	roster map[string]*Room
}

func NewRoomRegistry() *Registry {
	return &Registry{roster: make(map[string]*Room)}
}

func (registry *Registry) Add(room *Room) {
	registry.Lock()
	defer registry.Unlock()

	registry.roster[room.id] = room
}

func (registry *Registry) Find(id string) (*Room, bool) {
	registry.RLock()
	defer registry.RUnlock()

	room, ok := registry.roster[id]

	return room, ok
}
