package network

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/rhuantac/go-shooter/server"
)

type GameServer struct {
	SocketManager *SocketManager
	GameManager   *server.GameManager
}

func (gs GameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"localhost:5173"}})
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer c.CloseNow()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*210)
		defer cancel()
		var sAction SocketAction
		err = wsjson.Read(ctx, c, &sAction)
		if err != nil {
			log.Printf("failed to read JSON %v: %v", r.RemoteAddr, err)
			return
		}
		log.Printf("received: %v", sAction.Type)

		switch sAction.Type {
		case Init:
			var action InitAction
			json.Unmarshal(sAction.Action, &action)
			gs.SocketManager.InitPlayer(action, c)
			gs.GameManager.InitPlayer(action.Id, action.Name)
		case Movement, Rotate:
			var action server.Action
			json.Unmarshal(sAction.Action, &action)
			gs.GameManager.PerformAction(action)
		default:
			log.Printf("Unknown action received: %v", sAction.Type)

		}

	}
}
