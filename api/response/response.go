package response

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
	Id   string `json:"id"`
	Name string `json:"name"`
}
