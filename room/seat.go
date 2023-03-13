package room

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/violarium/poplan/user"
)

type Seat struct {
	mu                sync.RWMutex
	user              *user.User
	room              *Room
	voteIndex         int
	voted             bool
	activeSubscribers int32
}

func NewSeat(room *Room, u *user.User) *Seat {
	return &Seat{room: room, user: u}
}

func (s *Seat) User() *user.User {
	return s.user
}

func (s *Seat) PublicVote() Vote {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.room.Status() == StatusVoted {
		return s.room.voteTemplate.Votes[s.voteIndex]
	}

	return VoteUnknown
}

func (s *Seat) PrivateVote() Vote {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.room.voteTemplate.Votes[s.voteIndex]
}

func (s *Seat) SetVoteIndex(voteIndex int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if voteIndex >= len(s.room.voteTemplate.Votes) {
		return errors.New("invalid vote index")
	}

	s.voteIndex = voteIndex

	return nil
}

func (s *Seat) VoteIndex() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.voteIndex
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
