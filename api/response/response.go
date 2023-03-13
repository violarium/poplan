package response

import (
	"github.com/violarium/poplan/room"
	"github.com/violarium/poplan/user"
)

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
	Id            string      `json:"id"`
	Name          string      `json:"name"`
	Status        room.Status `json:"status"`
	Owner         bool        `json:"owner"`
	Seats         []Seat      `json:"seats"`
	TemplateTitle string      `json:"templateTitle"`
	VoteCards     []VoteCard  `json:"voteCards"`
}

type VoteCard struct {
	Vote   Vote `json:"vote"`
	Active bool `json:"active"`
}

type VoteTemplate struct {
	Title string `json:"title"`
	Votes []Vote `json:"votes"`
}

type VoteTemplateList struct {
	Templates []VoteTemplate `json:"templates"`
}

type Vote struct {
	Value float32 `json:"value"`
	Type  string  `json:"type"`
}

type Seat struct {
	User       User `json:"user"`
	Vote       Vote `json:"vote"`
	Voted      bool `json:"voted"`
	VoteOpened bool `json:"voteOpened"`
	Owner      bool `json:"owner"`
	Active     bool `json:"active"`
}

func NewRoom(r *room.Room, u *user.User) Room {
	roomResponse := Room{
		Id:            r.Id(),
		Name:          r.Name(),
		Status:        r.Status(),
		Owner:         r.Owner() == u,
		Seats:         BuildSeats(r, u),
		TemplateTitle: r.VoteTemplate().Title,
		VoteCards:     BuildVoteCards(r, u),
	}

	return roomResponse
}

func BuildSeats(r *room.Room, u *user.User) []Seat {
	roomStatus := r.Status()

	roomSeats := r.Seats()
	seats := make([]Seat, 0, len(roomSeats))

	for _, s := range roomSeats {
		seatUser := s.User()

		var seatVote room.Vote
		var voteOpened bool
		if seatUser == u {
			seatVote = s.PrivateVote()
			voteOpened = true
		} else {
			seatVote = s.PublicVote()
			voteOpened = roomStatus == room.StatusVoted
		}

		seats = append(seats, Seat{
			User: User{
				Id:   seatUser.Id(),
				Name: seatUser.Name(),
			},
			Vote:       NewVote(seatVote),
			Voted:      s.Voted(),
			VoteOpened: voteOpened,
			Owner:      seatUser == r.Owner(),
			Active:     s.Active(),
		})
	}

	return seats
}

func BuildVoteCards(r *room.Room, u *user.User) []VoteCard {
	uSeat, uSeatFound := r.ParticipantSeat(u)

	templateVotes := r.VoteTemplate().Votes
	voteCards := make([]VoteCard, 0, len(templateVotes))

	for i, v := range templateVotes {
		voteCards = append(voteCards, VoteCard{
			Vote:   NewVote(v),
			Active: uSeatFound && uSeat.Voted() && uSeat.VoteIndex() == i,
		})
	}

	return voteCards
}

func NewVoteTemplate(template room.VoteTemplate) VoteTemplate {
	votes := template.Votes
	voteResponses := make([]Vote, 0, len(votes))
	for _, v := range votes {
		voteResponses = append(voteResponses, NewVote(v))
	}

	return VoteTemplate{
		Title: template.Title,
		Votes: voteResponses,
	}
}

func NewVote(vote room.Vote) Vote {
	return Vote{
		Value: vote.Value(),
		Type:  vote.Type(),
	}
}
