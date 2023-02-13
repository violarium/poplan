package main

import (
	"sync"

	"github.com/google/uuid"
)

type RoomStatus uint

const (
	RoomStatusVoting RoomStatus = iota + 1
	RoomStatusVoted
)

type Seat struct {
	user  *User
	vote  uint
	voted bool
}

// room has status: voting, voted
// room includes users with their votes - seats
// actions: user vote, end vote (->voted), reset (->voting)

type Room struct {
	sync.Mutex
	id     string
	status RoomStatus
	owner  *User
	seats  []*Seat
}

func NewRoom(owner *User) *Room {
	seats := []*Seat{{user: owner}}
	room := &Room{
		id:     uuid.NewString(),
		status: RoomStatusVoting,
		owner:  owner,
		seats:  seats,
	}

	return room
}

func (room *Room) join(user *User) {
	room.Lock()
	defer room.Unlock()

	for _, s := range room.seats {
		if s.user == user {
			return
		}
	}

	seat := &Seat{user: user}
	room.seats = append(room.seats, seat)
}

func (room *Room) remove(user *User) {
	room.Lock()
	defer room.Unlock()

	newSeats := make([]*Seat, 0, cap(room.seats))
	for _, s := range room.seats {
		if s.user == user {
			continue
		}
		newSeats = append(newSeats, s)
	}

	room.seats = newSeats
}
