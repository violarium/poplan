package request

type Register struct {
	Name string `json:"name"`
}

type CreateRoom struct {
	Name         string `json:"name"`
	VoteTemplate int    `json:"voteTemplate"`
}

type UpdateRoom struct {
	Name string `json:"name"`
}

type Vote struct {
	Value int `json:"value"`
}
