package network

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/rhuantac/go-shooter/server"
)

type GameServer struct {
	Manager     *SocketManager
	GameManager *server.GameManager
}

func (gs GameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"localhost:5173"}})
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer c.CloseNow()
	gs.Manager.addConnection(c)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*210)
		defer cancel()
		var v server.Action
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			log.Printf("failed to read JSON %v: %v", r.RemoteAddr, err)
			return
		}
		log.Printf("received: %v", v)
		gs.GameManager.PerformAction(v)

	}
}

type SocketManager struct {
	WritePump   chan server.Snapshot
	connections []*websocket.Conn
}

func (sm *SocketManager) addConnection(c *websocket.Conn) {
	sm.connections = append(sm.connections, c)
}

func (sm *SocketManager) Start() {
	for {
		content := <-sm.WritePump
		for _, conn := range sm.connections {
			//log.Printf("Enviando content: %v para a conn: %v", content, conn)
			ctx := context.Background()
			wsjson.Write(ctx, conn, content)
		}
	}
}
