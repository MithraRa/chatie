package domain

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type PrivateRoom struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Clients [2]*Client `json:"clients"`
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
