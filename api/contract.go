package api

type Home struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Message struct {
	Message string `json:"message"`
}

type Register struct {
	Name string `json:"name"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Registration struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
