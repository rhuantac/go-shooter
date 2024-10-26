package server

import (
	"time"

	"github.com/ByteArena/box2d"
)

type PlayerColor string

const (
	Purple PlayerColor = "CC2B52"
	Yellow PlayerColor = "FABC3F"
	Green  PlayerColor = "B1D690"
	Blue   PlayerColor = "3E92CC"
)

var PlayerColors = []PlayerColor{Purple, Yellow, Green, Blue}

type UserData struct {
	Name  string
	Color PlayerColor
}
type Snapshot struct {
	Objects []State
}

type State struct {
	PosX, PosY float64
	Id         string
	Name       string
	Color      PlayerColor
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

func (gm *GameManager) InitPlayer(id string, name string) {
	color := PlayerColors[len(gm.playerStore)%len(PlayerColors)]
	p := gm.gameProcessor.CreateCharacter()
	data := UserData{Name: name, Color: color}
	p.SetUserData(data)
	gm.playerStore[id] = p
}

func (gm *GameManager) PerformAction(a Action) {
	gm.actionQueue.Push(a)
}

func (gm *GameManager) GetSnapshots() []Snapshot {
	return gm.snapshotQueue.PopAll()
}

func CreateNewGame(processor Processor) GameManager {
	store := make(PlayerStore, 0)
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
			err := worldProcessor.Process(playerStore, actions)
			if err != nil {
				panic(err)
			}
			snapshot := Snapshot{Objects: []State{}}
			for id, p := range playerStore {
				userData := p.GetUserData().(UserData)
				snapshot.Objects = append(snapshot.Objects, State{Id: id, Name: userData.Name, PosX: p.GetPosition().X, PosY: p.GetPosition().Y, Color: userData.Color})
			}
			snapshotQueue.Push(snapshot)
			accumulator -= dt
		}

	}
}
