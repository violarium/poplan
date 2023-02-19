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

type Room struct {
	mu     sync.RWMutex
	id     string
	name   string
	status Status
	owner  *user.User
	seats  []*Seat
}

func NewRoom(owner *user.User, name string) *Room {
	room := &Room{
		id:     uuid.NewString(),
		name:   name,
		status: StatusVoting,
		owner:  owner,
	}
	room.seats = append(room.seats, NewSeat(room, owner))

	return room
}

func (room *Room) Id() string {
	return room.id
}

func (room *Room) Name() string {
	room.mu.RLock()
	defer room.mu.RUnlock()

	return room.name
}

func (room *Room) Status() Status {
	room.mu.RLock()
	defer room.mu.RUnlock()

	return room.status
}

func (room *Room) SetName(name string) {
	room.mu.Lock()
	defer room.mu.Unlock()

	room.name = name
}

func (room *Room) Owner() *user.User {
	return room.owner
}

func (room *Room) HasParticipant(participant *user.User) bool {
	for _, seat := range room.seats {
		if seat.user == participant {
			return true
		}
	}

	return false
}

func (room *Room) Join(participant *user.User) {
	room.mu.Lock()
	defer room.mu.Unlock()

	for _, s := range room.seats {
		if s.user == participant {
			return
		}
	}

	room.seats = append(room.seats, NewSeat(room, participant))
}

func (room *Room) Leave(participant *user.User) {
	room.mu.Lock()
	defer room.mu.Unlock()

	newSeats := make([]*Seat, 0, cap(room.seats))
	for _, s := range room.seats {
		if s.user == participant {
			continue
		}
		newSeats = append(newSeats, s)
	}

	room.seats = newSeats
}

func (room *Room) Vote(participant *user.User, vote uint) {
	room.mu.Lock()
	defer room.mu.Unlock()

	if room.status != StatusVoting {
		return
	}

	seat, seatFound := room.getSeatFor(participant)
	if !seatFound {
		return
	}

	seat.SetVote(vote)
	seat.SetVoted(true)
}

func (room *Room) Seats() []*Seat {
	room.mu.RLock()
	defer room.mu.RUnlock()

	return room.seats
}

func (room *Room) EndVote() {
	room.mu.Lock()
	defer room.mu.Unlock()

	room.status = StatusVoted
}

func (room *Room) Reset() {
	room.mu.Lock()
	defer room.mu.Unlock()

	room.status = StatusVoting
	for _, s := range room.seats {
		s.SetVote(0)
		s.SetVoted(false)
	}
}

func (room *Room) getSeatFor(participant *user.User) (*Seat, bool) {
	for _, s := range room.seats {
		if s.user == participant {
			return s, true
		}
	}

	return nil, false
}
