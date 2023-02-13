package main

import (
	"sync"
)

// todo: register limit for rooms

type RoomRegistry struct {
	sync.RWMutex
	roster map[string]*Room
}

func NewRoomRegistry() *RoomRegistry {
	return &RoomRegistry{roster: make(map[string]*Room)}
}

func (registry *RoomRegistry) add(room *Room) {
	registry.Lock()
	defer registry.Unlock()

	registry.roster[room.id] = room
}

func (registry *RoomRegistry) find(id string) (*Room, bool) {
	registry.RLock()
	defer registry.RUnlock()

	room, ok := registry.roster[id]

	return room, ok
}
