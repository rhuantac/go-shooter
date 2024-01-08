package server

import (
	"time"

	"github.com/ByteArena/box2d"
)

type Snapshot struct {
	Objects []State
}

type State struct {
	PosX, PosY float64
	Name       string
}

type PlayerStore map[string]*box2d.B2Body
type GameManager struct {
	gameProcessor Processor
	actionQueue   *Queue[Action]
	snapshotQueue *Queue[Snapshot]
	playerStore   PlayerStore
}

func (gm *GameManager) Start() chan struct{} {
	quitChan := make(chan struct{})
	go gameLoop(gm.gameProcessor, gm.playerStore, gm.actionQueue, gm.snapshotQueue, quitChan)
	return quitChan
}

func (gm *GameManager) PerformAction(a Action) {
	gm.actionQueue.Push(a)
}

func (gm *GameManager) GetSnapshots() []Snapshot {
	return gm.snapshotQueue.PopAll()
}

func CreateNewGame(processor Processor) GameManager {
	store := make(PlayerStore, 0)
	p := processor.CreateCharacter()
	store["John"] = p
	gameActionQueue := NewQueue[Action]()
	snapshotQueue := NewQueue[Snapshot]()
	return GameManager{gameProcessor: processor, playerStore: store, actionQueue: gameActionQueue, snapshotQueue: snapshotQueue}
}

func gameLoop(processor Processor, players PlayerStore, actionQueue *Queue[Action], snapshotQueue *Queue[Snapshot], quit chan struct{}) {
	for {
		select {
		case <-quit:
			return
		default:
			queueProcess(processor, players, actionQueue, snapshotQueue)
		}

	}
}

func queueProcess(worldProcessor Processor, playerStore map[string]*box2d.B2Body, actionQueue *Queue[Action], snapshotQueue *Queue[Snapshot]) {
	dt := int64(16) // 1/60 in ms
	currentTime := time.Now().UnixMilli()
	accumulator := int64(0)
	for {
		newTime := time.Now().UnixMilli()
		iterationTime := newTime - currentTime
		
		currentTime = newTime

		accumulator += iterationTime

		for accumulator >= dt {
			actions := actionQueue.PopAll()
			worldProcessor.Process(playerStore, actions)
			snapshot := Snapshot{Objects: []State{}}
			for id, p := range playerStore {
				snapshot.Objects = append(snapshot.Objects, State{Name: id, PosX: p.GetPosition().X, PosY: p.GetPosition().Y})
			}
			snapshotQueue.Push(snapshot)
			accumulator -= dt
		}

	}
}
