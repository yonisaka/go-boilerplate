package ws

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yonisaka/go-boilerplate/pkg/logger"
)

type payloadMsgWrap struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type payloadMsgPeer struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
}

type payloadMsgChat struct {
	PeerID     string `json:"peer_id"`
	PeerHandle string `json:"peer_handle"`
	Msg        string `json:"message"`
}

// peerRequest represents a peer request (join, leave etc.) that's processed
// by a Room.
type peerRequest struct {
	reqType string
	peer    *Peer
}

type Room struct {
	ID string

	hub            *Hub
	broadcastQueue chan []byte
	peers          map[*Peer]bool

	peerQueue chan peerRequest

	disposeSig chan bool
	closed     bool

	// Message / payload cache.
	payloadCache [][]byte

	timestamp time.Time
}

func NewRoom(id string, h *Hub) *Room {
	return &Room{
		ID:             id,
		hub:            h,
		peers:          make(map[*Peer]bool, 100),
		broadcastQueue: make(chan []byte, 100),
		peerQueue:      make(chan peerRequest, 100),
		disposeSig:     make(chan bool),
		payloadCache:   make([][]byte, 0, h.cfg.Ws.MaxCachedMessages),
	}
}

// AddPeer adds a new peer to the room given a WS connection from an HTTP
// handler.
func (r *Room) AddPeer(id, handle string, ws *websocket.Conn) {
	r.queuePeerRequest(TypePeerJoin, newPeer(id, handle, ws, r))
}

// Dispose signals the room to notify all connected peer messages, and dispose
// of itself.
func (r *Room) Dispose() {
	r.disposeSig <- true
}

// Broadcast broadcasts a message to all connected peers.
func (r *Room) Broadcast(data []byte, record bool) {
	r.broadcastQueue <- data
	if record {
		r.recordMsgPayload(data)
	}
}

//nolint:funlen
func (r *Room) run() {
loop:
	for {
		select {
		// Dispose request.
		case <-r.disposeSig:
			break loop

		// Incoming peer request.
		case req, ok := <-r.peerQueue:
			if !ok {
				break loop
			}

			logger.Info(fmt.Sprintf("room %s received peer request %s", r.ID, req.reqType))

			switch req.reqType {
			// A new peer has joined.
			case TypePeerJoin:
				// Room's capacity is exchausted. Kick the peer out.
				if len(r.peers) >= r.hub.cfg.Ws.MaxPeerPerRoom {
					err := req.peer.writeWSControl(websocket.FormatCloseMessage(websocket.CloseNormalClosure, TypeRoomFull))
					if err != nil {
						logger.Error(fmt.Sprintf("write close message error: %s", err))
					}

					err = req.peer.ws.Close()
					if err != nil {
						logger.Error(fmt.Sprintf("close websocket error: %s", err))
					}

					continue
				}

				r.peers[req.peer] = true
				go req.peer.RunListener()
				go req.peer.RunWriter()

				// Send the peer its info.
				req.peer.SendData(r.makePeerUpdatePayload(req.peer, TypePeerInfo))

				// Send the peer last N message.
				if r.hub.cfg.Ws.MaxCachedMessages > 0 {
					for _, b := range r.payloadCache {
						req.peer.SendData(b)
					}
				}

				// Notify all peers of the new addition.
				r.Broadcast(r.makePeerUpdatePayload(req.peer, TypePeerJoin), true)

				logger.Info(fmt.Sprintf("%s joined %s", req.peer.Handle, r.ID))

			// A peer has left.
			case TypePeerLeave:
				r.removePeer(req.peer)
				r.Broadcast(r.makePeerUpdatePayload(req.peer, TypePeerLeave), true)
				logger.Info(fmt.Sprintf("%s left %s", req.peer.Handle, r.ID))

			// A peer has requested the room's peer list.
			case TypePeerList:
				req.peer.SendData(r.makePeerListPayload())
			}

		// Fanout broadcast to all peers.
		case m, ok := <-r.broadcastQueue:
			if !ok {
				break loop
			}

			for p := range r.peers {
				logger.Info(fmt.Sprintf("room %s broadcast message to %s", r.ID, p.Handle))
				p.SendData(m)
			}

			// Extend the room's expiry (once every 30 seconds).
			if time.Since(r.timestamp) > time.Duration(30)*time.Second {
				r.timestamp = time.Now()
				//r.extendTTL()
			}

		// Kill the room after the inactivity period.
		case <-time.After(time.Duration(r.hub.cfg.Ws.RoomAge) * time.Hour):
			logger.Error(fmt.Sprintf("room %s timed out", r.ID))
			break loop
		}
	}

	logger.Info(fmt.Sprintf("stopped room: %v", r.ID))
	r.remove()
}

// extendTTL extends a room's TTL in the store.
func (r *Room) extendTTL() {
	//r.hub.Store.ExtendRoomTTL(r.ID, r.hub.cfg.RoomAge)
}

// remove disposes a room by notifying and disconnecting all peers and
// removing the room from the store.
func (r *Room) remove() {
	r.closed = true

	// Close all peer WS connections.
	for peer := range r.peers {
		err := peer.writeWSControl(websocket.FormatCloseMessage(websocket.CloseNormalClosure, TypeRoomDispose))
		if err != nil {
			logger.Error(fmt.Sprintf("write close message error: %s", err))
		}

		delete(r.peers, peer)
	}

	// Close all room channels.
	close(r.broadcastQueue)
	close(r.peerQueue)
	r.hub.removeRoom(r.ID)
}

// recordMsgPayload records message payloads (events) sent out. It maintains last
// N messages to be sent to new users when they join.
func (r *Room) recordMsgPayload(b []byte) {
	if r.hub.cfg.Ws.MaxCachedMessages == 0 {
		return
	}

	n := len(r.payloadCache)
	if n >= r.hub.cfg.Ws.MaxCachedMessages {
		r.payloadCache = r.payloadCache[1:]
	}

	r.payloadCache = append(r.payloadCache, b)
}

// queuePeerRequest queues a peer addition / removal request to the room.
func (r *Room) queuePeerRequest(reqType string, p *Peer) {
	if r.closed {
		return
	}

	logger.Info(fmt.Sprintf("room %s queued peer request %s", r.ID, reqType))
	p.room.peerQueue <- peerRequest{reqType: reqType, peer: p}
}

// removePeer removes a peer from the room and broadcasts a message to the
// room notifying all peers of the action.
func (r *Room) removePeer(p *Peer) {
	close(p.dataQueue)
	delete(r.peers, p)
}

// sendPeerList sends the peer list to the given peer.
func (r *Room) sendPeerList(p *Peer) {
	r.peerQueue <- peerRequest{reqType: TypePeerList, peer: p}
}

// makePeerListPayload prepares a message payload with the list of peers.
func (r *Room) makePeerListPayload() []byte {
	peers := make([]payloadMsgPeer, 0, len(r.peers))
	for p := range r.peers {
		peers = append(peers, payloadMsgPeer{ID: p.ID, Handle: p.Handle})
	}

	return r.makePayload(peers, TypePeerList)
}

// makePeerUpdatePayload prepares a message payload representing a peer
// join / leave event.
func (r *Room) makePeerUpdatePayload(p *Peer, peerUpdateType string) []byte {
	d := payloadMsgPeer{
		ID:     p.ID,
		Handle: p.Handle,
	}

	return r.makePayload(d, peerUpdateType)
}

// makeMessagePayload prepares a chat message.
func (r *Room) makeMessagePayload(msg string, p *Peer) []byte {
	d := payloadMsgChat{
		PeerID:     p.ID,
		PeerHandle: p.Handle,
		Msg:        msg,
	}

	return r.makePayload(d, TypeMessage)
}

// makePayload prepares a message payload.
func (r *Room) makePayload(data interface{}, typ string) []byte {
	m := payloadMsgWrap{
		Timestamp: time.Now(),
		Type:      typ,
		Data:      data,
	}

	b, err := json.Marshal(m)
	if err != nil {
		logger.Error(err)
		return nil
	}

	return b
}
