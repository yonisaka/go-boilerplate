package ws

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/yonisaka/go-boilerplate/pkg/logger"
)

const (
	TypeTyping          = "typing"
	TypeMessage         = "message"
	TypePeerList        = "peer.list"
	TypePeerInfo        = "peer.info"
	TypePeerJoin        = "peer.join"
	TypePeerLeave       = "peer.leave"
	TypePeerRateLimited = "peer.ratelimited"
	TypeRoomDispose     = "room.dispose"
	TypeRoomFull        = "room.full"
	TypeNotice          = "notice"
	TypeHandle          = "handle"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("upgrade error: %s", err))
		return
	}

	peerID := r.Header.Get("X-Peer-ID")
	roomID := r.Header.Get("X-Room-ID")

	room := hub.InitRoom(roomID)
	sessionID := uuid.New().String()

	room.AddPeer(sessionID, peerID, ws)
}
