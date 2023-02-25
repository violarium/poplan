package room

import (
	"sync"

	"github.com/violarium/poplan/user"
)

type SeatSubscriber struct {
	Notifications chan bool
}

type Seat struct {
	mu            sync.RWMutex
	user          *user.User
	room          *Room
	vote          Vote
	voted         bool
	subscribers   map[*SeatSubscriber]bool
	subscribersMu sync.Mutex
}

func NewSeat(room *Room, u *user.User) *Seat {
	return &Seat{room: room, user: u, subscribers: make(map[*SeatSubscriber]bool)}
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

func (s *Seat) addSubscriber() *SeatSubscriber {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	subscriber := &SeatSubscriber{Notifications: make(chan bool, 256)}
	s.subscribers[subscriber] = true

	return subscriber
}

func (s *Seat) removeSubscriber(subscriber *SeatSubscriber) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	close(subscriber.Notifications)
	delete(s.subscribers, subscriber)
}

func (s *Seat) notifyAll() {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	for subscriber := range s.subscribers {
		select {
		case subscriber.Notifications <- true:
		default:
			close(subscriber.Notifications)
			delete(s.subscribers, subscriber)
		}
	}
}
