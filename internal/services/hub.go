package services

import (
	"chatie.com/internal/domain"
)

const (
	LeaveMessage = "user left the room"
)

type HubService struct {
	Hub *domain.Hub
}

func NewHubService(hub *domain.Hub) *HubService {
	return &HubService{
		Hub: hub,
	}
}

func (hs *HubService) Run() {
	for {
		select {
		case cl := <-hs.Hub.Register:
			if _, ok := hs.Hub.Rooms[cl.RoomID]; ok {
				room := hs.Hub.Rooms[cl.RoomID]
				if _, ok := room.Clients[cl.ID]; !ok {
					room.Clients[cl.ID] = cl
				}
			}
		case cl := <-hs.Hub.Unregister:
			if _, ok := hs.Hub.Rooms[cl.RoomID]; ok {
				if _, ok := hs.Hub.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					if len(hs.Hub.Rooms[cl.RoomID].Clients) != 0 {
						hs.Hub.Broadcast <- &domain.Message{
							Content:  LeaveMessage,
							RoomID:   cl.RoomID,
							Username: cl.Username,
						}
					}
					delete(hs.Hub.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
				}
			}
		case m := <-hs.Hub.Broadcast:
			if _, ok := hs.Hub.Rooms[m.RoomID]; ok {
				for _, cl := range hs.Hub.Rooms[m.RoomID].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
