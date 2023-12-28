package main

import (
	"time"

	"github.com/ByteArena/box2d"
)

type PlayerStore map[string]*box2d.B2Body
type GameManager struct {
	gameProcessor Processor
	playerStore   PlayerStore
}

func CreateNewGame(processor Processor) GameManager {
	store := make(PlayerStore, 0)
	store["John"] = processor.CreateCharacter()
	quitChan := make(chan struct{})
	go gameLoop(processor, store, quitChan, endGame)
	return GameManager{gameProcessor: processor, playerStore: store}
}

func gameLoop(processor Processor, players PlayerStore, quit chan struct{}, endGameFunc func()) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			for i := 0; i < 60; i++ {
				processor.Process(players, nil)
			}
		case <-quit:
			ticker.Stop()
			endGameFunc()
			return
		}
	}
}

func endGame() {
	println("End game called")
}
