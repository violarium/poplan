package room

import (
	"errors"
	"sync"

	"github.com/violarium/poplan/user"
	"github.com/violarium/poplan/util"
)

type Status uint

const (
	StatusVoting Status = iota + 1
	StatusVoted
)

type Room struct {
	mu           sync.RWMutex
	id           string
	name         string
	status       Status
	owner        *user.User
	seats        []*Seat
	voteTemplate VoteTemplate
}

func NewRoom(owner *user.User, name string, voteTemplate VoteTemplate) *Room {
	room := &Room{
		id:           util.GeneratePrettyId(8),
		name:         name,
		status:       StatusVoting,
		owner:        owner,
		voteTemplate: voteTemplate,
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

func (room *Room) VoteTemplate() VoteTemplate {
	return room.voteTemplate
}

func (room *Room) Status() Status {
	room.mu.RLock()
	defer room.mu.RUnlock()

	return room.status
}

func (room *Room) SetName(name string) {
	room.mu.Lock()
	defer room.mu.Unlock()

	defer room.notifyAll()

	room.name = name
}

func (room *Room) Owner() *user.User {
	return room.owner
}

func (room *Room) HasParticipant(participant *user.User) bool {
	room.mu.RLock()
	defer room.mu.RUnlock()

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

	defer room.notifyAll()

	for _, s := range room.seats {
		if s.user == participant {
			return
		}
	}

	room.seats = append(room.seats, NewSeat(room, participant))
}

func (room *Room) Subscribe(participant *user.User) (*SeatSubscriber, error) {
	room.mu.RLock()
	defer room.mu.RUnlock()

	defer room.notifyAll()

	seat, seatFound := room.getSeatFor(participant)
	if !seatFound {
		return nil, errors.New("user is not a participant")
	}

	return seat.addSubscriber(), nil
}

func (room *Room) Unsubscribe(subscriber *SeatSubscriber) {
	room.mu.RLock()
	defer room.mu.RUnlock()

	defer room.notifyAll()

	for _, s := range room.seats {
		s.removeSubscriber(subscriber)
	}
}

func (room *Room) Leave(participant *user.User) {
	room.mu.Lock()
	defer room.mu.Unlock()

	defer room.notifyAll()

	newSeats := make([]*Seat, 0, cap(room.seats))
	for _, s := range room.seats {
		if s.user == participant {
			continue
		}
		newSeats = append(newSeats, s)
	}

	room.seats = newSeats
}

func (room *Room) Vote(participant *user.User, voteIndex int) {
	room.mu.RLock()
	defer room.mu.RUnlock()

	defer room.notifyAll()

	if room.status != StatusVoting {
		return
	}

	seat, seatFound := room.getSeatFor(participant)
	if !seatFound {
		return
	}

	if voteIndex >= len(room.voteTemplate.Votes) {
		return
	}
	vote := room.voteTemplate.Votes[voteIndex]

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

	defer room.notifyAll()

	room.status = StatusVoted
}

func (room *Room) Reset() {
	room.mu.Lock()
	defer room.mu.Unlock()

	defer room.notifyAll()

	room.status = StatusVoting
	for _, s := range room.seats {
		s.SetVote(VoteUnknown)
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

func (room *Room) notifyAll() {
	for _, s := range room.seats {
		s.notifyAll()
	}
}
