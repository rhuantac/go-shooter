package network

import (
	"context"
	"encoding/json"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/rhuantac/go-shooter/server"
)

type ActionType int

const (
	Init     ActionType = 1
	Movement ActionType = 2
)

type PlayerConn struct {
	Id   string
	Conn *websocket.Conn
}

// Main packet structure. Params will contain the according packet type
type SocketAction struct {
	Type   ActionType
	Action json.RawMessage
}

type InitAction struct {
	Id   string
	Name string
}

type SocketManager struct {
	WritePump         chan server.Snapshot
	playerConnections map[string]PlayerConn
}

func (sm *SocketManager) addConnection(c PlayerConn) {
	if _, exists := sm.playerConnections[c.Id]; exists {
		return
	}
	sm.playerConnections[c.Id] = c
}

func (sm *SocketManager) InitPlayer(action InitAction, c *websocket.Conn) {
	sm.addConnection(PlayerConn{Id: action.Id, Conn: c})
}
func (sm *SocketManager) Start() {
	sm.playerConnections = make(map[string]PlayerConn)
	for {
		content := <-sm.WritePump
		for _, pConn := range sm.playerConnections {
			ctx := context.Background()
			wsjson.Write(ctx, pConn.Conn, content)
		}
	}
}
