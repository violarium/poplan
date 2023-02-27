package room

import (
	"sync"
	"sync/atomic"

	"github.com/violarium/poplan/user"
)

type Seat struct {
	mu                sync.RWMutex
	user              *user.User
	room              *Room
	vote              Vote
	voted             bool
	activeSubscribers int32
}

func NewSeat(room *Room, u *user.User) *Seat {
	return &Seat{room: room, user: u}
}

func (s *Seat) User() *user.User {
	return s.user
}

func (s *Seat) SecretVote() Vote {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.room.Status() == StatusVoted {
		return s.vote
	}

	return VoteUnknown
}

func (s *Seat) SetVote(vote Vote) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.vote = vote
}

func (s *Seat) Voted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.voted
}

func (s *Seat) SetVoted(voted bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.voted = voted
}

func (s *Seat) Active() bool {
	return atomic.LoadInt32(&s.activeSubscribers) > 0
}

func (s *Seat) IncActive() {
	atomic.AddInt32(&s.activeSubscribers, 1)
}

func (s *Seat) DecActive() {
	atomic.AddInt32(&s.activeSubscribers, -1)
}
