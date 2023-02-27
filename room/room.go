package room

import (
	"sync"

	"github.com/violarium/poplan/user"
	"github.com/violarium/poplan/util"
)

type Status uint

const (
	StatusVoting Status = iota + 1
	StatusVoted
)

type ChangeHandler struct {
	callback func(*Room)
}

func NewChangeHandler(callback func(*Room)) *ChangeHandler {
	return &ChangeHandler{callback: callback}
}

type Room struct {
	mu             sync.RWMutex
	id             string
	name           string
	status         Status
	owner          *user.User
	seats          []*Seat
	voteTemplate   VoteTemplate
	changeHandlers map[*ChangeHandler]bool
}

func NewRoom(owner *user.User, name string, voteTemplate VoteTemplate) *Room {
	room := &Room{
		id:             util.GeneratePrettyId(8),
		name:           name,
		status:         StatusVoting,
		owner:          owner,
		voteTemplate:   voteTemplate,
		changeHandlers: make(map[*ChangeHandler]bool),
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

	defer room.callChangeHandlers()

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

	defer room.callChangeHandlers()

	for _, s := range room.seats {
		if s.user == participant {
			return
		}
	}

	room.seats = append(room.seats, NewSeat(room, participant))
}

func (room *Room) IncActiveParticipant(participant *user.User) {
	room.mu.RLock()
	defer room.mu.RUnlock()

	seat, seatFound := room.getSeatFor(participant)
	if !seatFound {
		return
	}

	defer room.callChangeHandlers()
	seat.IncActive()
}

func (room *Room) DecActiveParticipant(participant *user.User) {
	room.mu.RLock()
	defer room.mu.RUnlock()

	seat, seatFound := room.getSeatFor(participant)
	if !seatFound {
		return
	}

	defer room.callChangeHandlers()
	seat.DecActive()
}

func (room *Room) Leave(participant *user.User) {
	room.mu.Lock()
	defer room.mu.Unlock()

	defer room.callChangeHandlers()

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

	defer room.callChangeHandlers()

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

	defer room.callChangeHandlers()

	room.status = StatusVoted
}

func (room *Room) Reset() {
	room.mu.Lock()
	defer room.mu.Unlock()

	defer room.callChangeHandlers()

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

func (room *Room) AddChangeHandler(changeHandler *ChangeHandler) {
	room.mu.Lock()
	defer room.mu.Unlock()

	room.changeHandlers[changeHandler] = true
}

func (room *Room) RemoveChangeHandler(changeHandler *ChangeHandler) {
	room.mu.Lock()
	defer room.mu.Unlock()

	delete(room.changeHandlers, changeHandler)
}

func (room *Room) callChangeHandlers() {
	for h := range room.changeHandlers {
		go h.callback(room)
	}
}
