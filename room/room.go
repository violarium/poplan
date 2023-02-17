package room

import (
	"sync"

	"github.com/google/uuid"
	"github.com/violarium/poplan/user"
)

type Status uint

const (
	StatusVoting Status = iota + 1
	StatusVoted
)

type Seat struct {
	user  *user.User
	vote  uint
	voted bool
}

type Room struct {
	sync.Mutex
	id     string
	name   string
	status Status
	owner  *user.User
	seats  []*Seat
}

func NewRoom(owner *user.User, name string) *Room {
	seats := []*Seat{{user: owner}}
	room := &Room{
		id:     uuid.NewString(),
		name:   name,
		status: StatusVoting,
		owner:  owner,
		seats:  seats,
	}

	return room
}

func (room *Room) Id() string {
	return room.id
}

func (room *Room) Name() string {
	return room.name
}

func (room *Room) Join(participant *user.User) {
	room.Lock()
	defer room.Unlock()

	for _, s := range room.seats {
		if s.user == participant {
			return
		}
	}

	seat := &Seat{user: participant}
	room.seats = append(room.seats, seat)
}

func (room *Room) Leave(participant *user.User) {
	room.Lock()
	defer room.Unlock()

	newSeats := make([]*Seat, 0, cap(room.seats))
	for _, s := range room.seats {
		if s.user == participant {
			continue
		}
		newSeats = append(newSeats, s)
	}

	room.seats = newSeats
}

func (room *Room) EndVote() {
	room.Lock()
	defer room.Unlock()

	room.status = StatusVoted
}

func (room *Room) Reset() {
	room.Lock()
	defer room.Unlock()

	room.status = StatusVoting
	for _, s := range room.seats {
		s.vote = 0
		s.voted = false
	}
}
