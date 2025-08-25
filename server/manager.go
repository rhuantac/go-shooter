package server

import (
	"fmt"
	"math/rand"
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
	Players []State
	Objects []State
}

type State struct {
	PosX, PosY float64
	Id         string
	Name       string
	Color      PlayerColor
	Rotation   float64
}

type PlayerStore map[string]*box2d.B2Body
type GameManager struct {
	gameProcessor Processor
	actionQueue   *Queue[Action]
	snapshotQueue *Queue[Snapshot]
	playerStore   PlayerStore
	objectStore   map[string]*box2d.B2Body
}

func (gm *GameManager) Start() chan struct{} {
	quitChan := make(chan struct{})
	go gameLoop(gm, quitChan)
	return quitChan
}

func (gm *GameManager) InitPlayer(id string, name string) {
	color := PlayerColors[len(gm.playerStore)%len(PlayerColors)]
	p := gm.gameProcessor.CreateCharacter()
	data := UserData{Name: name, Color: color}
	characterData := BodyData{ID: fmt.Sprintf("player-%d", rand.Int()), Type: "player", UserData: data}
	p.SetUserData(characterData)
	gm.playerStore[id] = p
}

func (gm *GameManager) PerformAction(a Action) {
	gm.actionQueue.Push(a)
}

func (gm *GameManager) GetSnapshots() []Snapshot {
	return gm.snapshotQueue.PopAll()
}

func CreateNewGame(processor Processor) GameManager {
	playerStore := make(PlayerStore, 0)
	objectStore := make(map[string]*box2d.B2Body, 0)
	gameActionQueue := NewQueue[Action]()
	snapshotQueue := NewQueue[Snapshot]()
	return GameManager{gameProcessor: processor, playerStore: playerStore, objectStore: objectStore, actionQueue: gameActionQueue, snapshotQueue: snapshotQueue}
}

func gameLoop(gameManager *GameManager, quit chan struct{}) {
	for {
		select {
		case <-quit:
			return
		default:
			queueProcess(gameManager)
		}

	}
}

func queueProcess(gameManager *GameManager) {
	dt := int64(33) // 1/60 in ms
	currentTime := time.Now().UnixMilli()
	accumulator := int64(0)
	for {
		newTime := time.Now().UnixMilli()
		iterationTime := newTime - currentTime

		currentTime = newTime

		accumulator += iterationTime

		for accumulator >= dt {
			actions := gameManager.actionQueue.PopAll()
			err := gameManager.gameProcessor.Process(gameManager.playerStore, gameManager.objectStore, actions)
			if err != nil {
				panic(err)
			}
			snapshot := Snapshot{Objects: []State{}, Players: []State{}}
			for id, p := range gameManager.playerStore {
				userData := p.GetUserData().(BodyData).UserData.(UserData)
				snapshot.Players = append(snapshot.Players, State{Id: id, Name: userData.Name, PosX: p.GetPosition().X, PosY: p.GetPosition().Y, Color: userData.Color, Rotation: p.GetAngle()})
			}
			for id, o := range gameManager.objectStore {
				snapshot.Objects = append(snapshot.Objects, State{Id: id, PosX: o.GetPosition().X, PosY: o.GetPosition().Y, Rotation: o.GetAngle(), Color: "FFFFFF"})
			}
			gameManager.snapshotQueue.Push(snapshot)
			accumulator -= dt
		}

	}
}
