package request

type Register struct {
	Name string `json:"name"`
}

type CreateRoom struct {
	Name string `json:"name"`
}
