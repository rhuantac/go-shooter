package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rhuantac/go-shooter/network"
	"github.com/rhuantac/go-shooter/server"
)

func main() {
	log.SetFlags(0)
	writePump := make(chan server.Snapshot)
	gameManager := runGame(writePump)
	err := runWs(writePump, gameManager)
	if err != nil {
		log.Fatal(err)
	}
}

// runWs starts a http.Server for the passed in address
// with all requests handled by echoServer.
func runWs(writePump chan server.Snapshot, gameManager *server.GameManager) error {

	addr := "localhost:5000"
	if len(os.Args) >= 2 {
		addr = os.Args[1]
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("listening on ws://%v", l.Addr())

	socketManager := network.SocketManager{
		WritePump: writePump,
	}
	go socketManager.Start()

	s := &http.Server{
		Handler:      network.GameServer{SocketManager: &socketManager, GameManager: gameManager},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}

func runGame(writePump chan<- server.Snapshot) *server.GameManager {
	processor := server.NewProcessor()
	manager := server.CreateNewGame(&processor)
	manager.Start()

	go func() {
		for {
			snaps := manager.GetSnapshots()
			if len(snaps) > 0 {
				for _, snap := range snaps {
					writePump <- snap
				}
			}
		}

	}()
	return &manager
}
