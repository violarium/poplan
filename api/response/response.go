package response

import "github.com/violarium/poplan/room"

type Home struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Message struct {
	Message string `json:"message"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Registration struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type Room struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Status       room.Status  `json:"status"`
	Seats        []Seat       `json:"seats"`
	VoteTemplate VoteTemplate `json:"voteTemplate"`
}

type VoteTemplate struct {
	Title string `json:"title"`
	Votes []Vote `json:"votes"`
}

type Vote struct {
	Value float32 `json:"value"`
	Type  string  `json:"type"`
}

type Seat struct {
	User  User `json:"user"`
	Vote  Vote `json:"vote"`
	Voted bool `json:"voted"`
	Owner bool `json:"owner"`
}

func NewRoom(r *room.Room) Room {
	seats := r.Seats()
	seatsResponse := make([]Seat, 0, len(seats))
	for _, s := range seats {
		seatsResponse = append(seatsResponse, Seat{
			User: User{
				Id:   s.User().Id(),
				Name: s.User().Name(),
			},
			Vote: Vote{
				Value: s.SecretVote().Value(),
				Type:  s.SecretVote().Type(),
			},
			Voted: s.Voted(),
			Owner: s.User() == r.Owner(),
		})
	}

	votes := r.VoteTemplate().Votes
	voteResponses := make([]Vote, 0, len(votes))
	for _, v := range votes {
		voteResponses = append(voteResponses, Vote{
			Value: v.Value(),
			Type:  v.Type(),
		})
	}

	roomResponse := Room{
		Id:     r.Id(),
		Name:   r.Name(),
		Status: r.Status(),
		Seats:  seatsResponse,
		VoteTemplate: VoteTemplate{
			Title: r.VoteTemplate().Title,
			Votes: voteResponses,
		},
	}

	return roomResponse
}
