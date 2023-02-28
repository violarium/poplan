package room

type Subscriber struct {
	notifications chan bool
}

func NewSubscriber() *Subscriber {
	return &Subscriber{notifications: make(chan bool)}
}

func (s *Subscriber) Notifications() <-chan bool {
	return s.notifications
}
