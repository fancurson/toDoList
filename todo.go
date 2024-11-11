package todo

type TodoList struct {
	Id          int    `json:"-"`
	Title       string `json: "title"`
	Description string `json: "description"`
}

type UserList struct {
	Id     int
	UserId string
	ListId string
}

type ListsItem struct {
	Id       int
	ListenId string
	ItemId   string
}
