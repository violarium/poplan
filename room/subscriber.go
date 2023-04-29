package room

type Subscriber struct {
	notifications chan bool
	onBlocked     func()
}

func NewSubscriber(bufferSize int, onBlocked func()) *Subscriber {
	return &Subscriber{notifications: make(chan bool, bufferSize), onBlocked: onBlocked}
}

func (s *Subscriber) Notifications() <-chan bool {
	return s.notifications
}
