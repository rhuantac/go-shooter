package main

import (
	"time"

	"github.com/ByteArena/box2d"
)

type Snapshot struct {
}
type PlayerStore map[string]*box2d.B2Body
type GameManager struct {
	gameProcessor Processor
	actionQueue   *ActionQueue
	playerStore   PlayerStore
}

func CreateNewGame(processor Processor) GameManager {
	store := make(PlayerStore, 0)
	store["John"] = processor.CreateCharacter()
	quitChan := make(chan struct{})
	gameActionQueue := &ActionQueue{}
	go gameLoop(processor, store, gameActionQueue, quitChan, endGame)
	return GameManager{gameProcessor: processor, playerStore: store, actionQueue: gameActionQueue}
}

func gameLoop(processor Processor, players PlayerStore, actionQueue *ActionQueue, quit chan struct{}, endGameFunc func()) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			queueProcess(processor, players, actionQueue, &SnapshotQueue{})
		case <-quit:
			ticker.Stop()
			endGameFunc()
			return
		}
	}
}

func queueProcess(worldProcessor Processor, playerStore map[string]*box2d.B2Body, actionQueue *ActionQueue, snapshotQueue *SnapshotQueue) {
	for i := 0; i < 60; i++ {
		actions := actionQueue.PopAll()
		worldProcessor.Process(playerStore, actions)
		snapshotQueue.Push(Snapshot{})
	}
}

func endGame() {
	println("End game called")
}
