package ws

import (
	"sync"

	"github.com/yonisaka/go-boilerplate/config"
)

type Hub struct {
	rooms map[string]*Room
	cfg   *config.Config

	mut sync.RWMutex
}

func NewHub(cfg *config.Config) *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
		cfg:   cfg,
	}
}

// InitRoom returns a room with the given id.
func (h *Hub) InitRoom(id string) *Room {
	h.mut.RLock()
	room, ok := h.rooms[id]
	h.mut.RUnlock()

	if ok {
		return room
	}

	room = NewRoom(id, h)
	h.mut.Lock()
	h.rooms[id] = room
	h.mut.Unlock()

	go room.run()

	return room
}

func (h *Hub) getRooms() []*Room {
	h.mut.RLock()
	out := make([]*Room, 0, len(h.rooms))

	for _, r := range h.rooms {
		out = append(out, r)
	}
	h.mut.RUnlock()

	return out
}

// removeRoom removes a room from the hub and the store.
func (h *Hub) removeRoom(id string) {
	h.mut.Lock()
	delete(h.rooms, id)
	h.mut.Unlock()
}
