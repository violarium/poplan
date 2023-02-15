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

// room has status: voting, voted
// room includes users with their votes - seats
// actions: user vote, end vote (->voted), reset (->voting)

type Room struct {
	sync.Mutex
	id     string
	status Status
	owner  *user.User
	seats  []*Seat
}

func NewRoom(owner *user.User) *Room {
	seats := []*Seat{{user: owner}}
	room := &Room{
		id:     uuid.NewString(),
		status: StatusVoting,
		owner:  owner,
		seats:  seats,
	}

	return room
}

func (room *Room) join(participant *user.User) {
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

func (room *Room) remove(participant *user.User) {
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
