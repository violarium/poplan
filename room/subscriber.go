package room

type Subscriber struct {
	notifications chan bool
}

func (s *Subscriber) Notifications() <-chan bool {
	return s.notifications
}
