package request

type Register struct {
	Name string `json:"name"`
}

type CreateRoom struct {
	Name string `json:"name"`
}

type UpdateRoom struct {
	Name string `json:"name"`
}

type Vote struct {
	Value uint `json:"value"`
}
