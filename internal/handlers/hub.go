package handlers

import (
	"encoding/json"
	"net/http"

	"chatie.com/internal/domain"
	"chatie.com/internal/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
)

const (
	UserJoinedTheChat = "a user joined the room"
)

type HubHandler struct {
	hub *services.HubService
}

func NewHubHandler(h *services.HubService) *HubHandler {
	return &HubHandler{
		h,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *HubHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var roomReq domain.CreateRoomReq
	if err := json.NewDecoder(r.Body).Decode(&roomReq); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err)
		return
	}

	h.hub.Hub.Rooms[roomReq.ID] = &domain.Room{
		ID:      roomReq.ID,
		Name:    roomReq.Name,
		Clients: make(map[string]*domain.Client),
	}

	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, roomReq)
}

func (h *HubHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err)
		return
	}

	roomID := chi.URLParam(r, "roomId")
	clientID := chi.URLParam(r, "userId")
	username := chi.URLParam(r, "username")

	cl := &domain.Client{
		Conn:     conn,
		Message:  make(chan *domain.Message),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &domain.Message{
		Content:  UserJoinedTheChat,
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Hub.Register <- cl
	h.hub.Hub.Broadcast <- m

	go cl.WriteMessage()
	cl.ReadMessage(h.hub.Hub)
}

func (h *HubHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := make([]domain.RoomRes, 0)

	for _, r := range h.hub.Hub.Rooms {
		rooms = append(rooms, domain.RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, rooms)
}

func (h *HubHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	var clients []domain.ClientRes
	roomId := chi.URLParam(r, "roomId")

	if _, ok := h.hub.Hub.Rooms[roomId]; !ok {
		clients = make([]domain.ClientRes, 0)
		render.Status(r, http.StatusOK)
		render.JSON(w, r, clients)
	}

	for _, c := range h.hub.Hub.Rooms[roomId].Clients {
		clients = append(clients, domain.ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, clients)
}
