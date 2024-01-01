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
	actionQueue   *Queue[Action]
	snapshotQueue *Queue[Snapshot]
	playerStore   PlayerStore
}

func (gm *GameManager) Start() chan struct{}{
	quitChan := make(chan struct{})
	go gameLoop(gm.gameProcessor, gm.playerStore, gm.actionQueue, gm.snapshotQueue, quitChan)
	return quitChan
}

func CreateNewGame(processor Processor) GameManager {
	store := make(PlayerStore, 0)
	store["John"] = processor.CreateCharacter()
	gameActionQueue := NewQueue[Action]()
	snapshotQueue := NewQueue[Snapshot]()
	return GameManager{gameProcessor: processor, playerStore: store, actionQueue: gameActionQueue, snapshotQueue: snapshotQueue}
}

func gameLoop(processor Processor, players PlayerStore, actionQueue *Queue[Action], snapshotQueue *Queue[Snapshot], quit chan struct{}) {
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

func queueProcess(worldProcessor Processor, playerStore map[string]*box2d.B2Body, actionQueue *Queue[Action], snapshotQueue *Queue[Snapshot]) {
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
