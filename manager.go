package main

import (
	"time"

	"github.com/ByteArena/box2d"
)

type Snapshot struct {
	state []State
}

type State struct {
	posX, posY float64
	name       string
}
type PlayerStore map[string]*box2d.B2Body
type GameManager struct {
	gameProcessor Processor
	actionQueue   *ActionQueue
	snapshotQueue *SnapshotQueue
	playerStore   PlayerStore
}

func CreateNewGame(processor Processor) GameManager {
	store := make(PlayerStore, 0)
	store["John"] = processor.CreateCharacter()
	quitChan := make(chan struct{})
	gameActionQueue := NewActionQueue()
	snapshotQueue := &SnapshotQueue{}
	go gameLoop(processor, store, gameActionQueue, snapshotQueue, quitChan)
	return GameManager{gameProcessor: processor, playerStore: store, actionQueue: gameActionQueue, snapshotQueue: snapshotQueue}
}

func gameLoop(processor Processor, players PlayerStore, actionQueue *ActionQueue, snapshotQueue *SnapshotQueue, quit chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			queueProcess(processor, players, actionQueue, snapshotQueue)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func queueProcess(worldProcessor Processor, playerStore map[string]*box2d.B2Body, actionQueue *ActionQueue, snapshotQueue *SnapshotQueue) {
	for i := 0; i < 60; i++ {
		actions := actionQueue.PopAll()
		worldProcessor.Process(playerStore, actions)
		snapshot := Snapshot{state: []State{}}
		for id, p := range playerStore {
			snapshot.state = append(snapshot.state, State{name: id, posX: p.GetPosition().X, posY: p.GetPosition().Y})
		}
		snapshotQueue.Push(snapshot)

	}
}
