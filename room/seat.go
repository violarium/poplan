package room

import (
	"errors"
	"sync"

	"github.com/violarium/poplan/user"
)

type Seat struct {
	mu            sync.RWMutex
	user          *user.User
	room          *Room
	voteIndex     int
	voted         bool
	subscribers   map[*Subscriber]bool
	subscribersMu sync.RWMutex
}

func NewSeat(room *Room, u *user.User) *Seat {
	return &Seat{
		room:        room,
		user:        u,
		subscribers: make(map[*Subscriber]bool),
	}
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
	s.subscribersMu.RLock()
	defer s.subscribersMu.RUnlock()

	return len(s.subscribers) > 0
}

func (s *Seat) Subscribe(subscriber *Subscriber) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	s.subscribers[subscriber] = true
}

func (s *Seat) Unsubscribe(subscriber *Subscriber) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	if _, ok := s.subscribers[subscriber]; ok {
		close(subscriber.notifications)
		delete(s.subscribers, subscriber)
	}
}

func (s *Seat) UnsubscribeAll() {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	for subscriber := range s.subscribers {
		close(subscriber.notifications)
	}
	s.subscribers = make(map[*Subscriber]bool)
}

func (s *Seat) notifySubscribers() {
	s.subscribersMu.RLock()
	defer s.subscribersMu.RUnlock()

	for subscriber := range s.subscribers {
		subscriber.notifications <- true
	}
}
