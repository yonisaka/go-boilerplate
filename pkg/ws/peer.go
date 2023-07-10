package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yonisaka/go-boilerplate/pkg/logger"
)

type Peer struct {
	ID     string
	Handle string

	ws *websocket.Conn

	// channel for outbound messages.
	dataQueue chan []byte

	// room that the peer belongs to.
	room *Room

	// Rate limiting.
	numMessages int
	lastMessage time.Time
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func newPeer(id, handle string, ws *websocket.Conn, room *Room) *Peer {
	return &Peer{
		ID:        id,
		Handle:    handle,
		ws:        ws,
		dataQueue: make(chan []byte, 100),
		room:      room,
	}
}

// RunListener is a blocking function that reads incoming messages from a peer's
// WS connection until its dropped or there's an error. This should be invoked
// as a goroutine.
func (p *Peer) RunListener() {
	defer func() {
		// WS connection is closed.
		logger.Info(fmt.Sprintf("Peer %s listener is closed", p.Handle))

		err := p.ws.Close()
		if err != nil {
			return
		}

		p.room.queuePeerRequest(TypePeerLeave, p)
	}()

	p.ws.SetReadLimit(int64(p.room.hub.cfg.Ws.MaxMessageLength))

	for {
		_, message, err := p.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error(fmt.Sprintf("Peer %s read error: %s", p.Handle, err))
			}

			break
		}

		message = bytes.TrimSpace(
			bytes.Replace(message, newline, space, -1),
		)

		p.processMessage(message)
	}
}

// RunWriter is a blocking function that writes messages in a peer's queue to the
// peer's WS connection. This should be invoked as a goroutine.
func (p *Peer) RunWriter() {
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			return
		}
	}(p.ws)

	for {
		select {
		// Wait for outgoing message to appear in the channel.
		case message, ok := <-p.dataQueue:
			if !ok {
				err := p.writeWSData(websocket.CloseMessage, []byte{})
				if err != nil {
					logger.Warn(fmt.Sprintf("Peer %s write error: %s", p.Handle, err))
					break
				}

				break
			}

			err := p.writeWSData(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}

// SendData queues a message to be written to the peer's WS.
func (p *Peer) SendData(b []byte) {
	p.dataQueue <- b
}

// writeWSData writes the given payload to the peer's WS connection.
func (p *Peer) writeWSData(msgType int, payload []byte) error {
	err := p.ws.SetWriteDeadline(time.Now().Add(time.Duration(p.room.hub.cfg.Ws.Timeout) * time.Second))
	if err != nil {
		return err
	}

	return p.ws.WriteMessage(msgType, payload)
}

// writeWSControl writes the given control payload to the peer's WS connection.
func (p *Peer) writeWSControl(payload []byte) error {
	return p.ws.WriteControl(websocket.CloseMessage, payload, time.Time{})
}

// processMessage processes incoming messages from peers.
func (p *Peer) processMessage(b []byte) {
	var m payloadMsgWrap

	if err := json.Unmarshal(b, &m); err != nil {
		// TODO: Respond
		return
	}

	switch m.Type {
	// Message to the room.
	case TypeMessage:
		// Check rate limits and update counters.
		now := time.Now()

		if p.numMessages > 0 { //nolint:nestif
			if (p.numMessages%p.room.hub.cfg.Ws.RateLimitMessages+1) >= p.room.hub.cfg.Ws.RateLimitMessages && time.Since(p.lastMessage) < time.Duration(p.room.hub.cfg.Ws.RateLimitInterval)*time.Second {

				err := p.writeWSControl(websocket.FormatCloseMessage(websocket.CloseNormalClosure, TypePeerRateLimited))
				if err != nil {
					logger.Warn(fmt.Sprintf("Peer %s write error: %s", p.Handle, err))
				}

				err = p.ws.Close()
				if err != nil {
					return
				}

				return
			}
		}

		p.lastMessage = now
		p.numMessages++

		msg, ok := m.Data.(string)
		if !ok {
			// TODO: Respond
			return
		}

		logger.Info(fmt.Sprintf("Process message: %s", msg))

		p.room.Broadcast(p.room.makeMessagePayload(msg, p), true)

	// "Typing" status.
	case TypeTyping:
		p.room.Broadcast(p.room.makePeerUpdatePayload(p, TypeTyping), false)

	// Request for peers list
	case TypePeerList:
		p.room.sendPeerList(p)

	// Dipose of a room.
	case TypeRoomDispose:
		p.room.Dispose()
	default:
	}
}
