package main

import (
	"errors"
	"image/color"
	"log"

	"github.com/rhuantac/go-shooter/server"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	screenWidth  = 1000
	screenHeight = 600
)

type Player struct {
	X, Y   float32
	states []server.State
}

func (p *Player) AddSnapshot(s server.State) {
	if len(p.states) >= 2 {
		// keep only the last 2 snaps
		p.states = append(p.states[:2], s)
		return
	}

	p.states = append(p.states, s)
}

func (p *Player) RemoveSnapshot() (server.State, error) {
	if len(p.states) == 0 {
		return server.State{}, errors.New("array is empty")
	}
	s := p.states[0]
	p.states = p.states[1:]
	return s, nil
}

func (p *Player) NextState() {
	next, err := p.RemoveSnapshot()
	if err != nil {
		return
	}
	p.X = float32(next.PosX)
	p.Y = float32(-1 * next.PosY) //Server coordinates are cartesian
}

type KeyController struct {
	keyMap map[server.InputType]bool
}

func (k *KeyController) keyPress(key server.InputType) {
	k.keyMap[key] = true
}

func (k *KeyController) keyRelease(key server.InputType) {
	k.keyMap[key] = false
}

func (k *KeyController) wasPressing(key server.InputType) bool {
	return k.keyMap[key]
}

type Game struct {
	manager *server.GameManager
	players map[string]*Player
	keys    KeyController
}

func (g *Game) Update() error {
	snaps := g.manager.GetSnapshots()
	for _, snap := range snaps {
		for _, state := range snap.Objects {
			if _, exists := g.players[state.Name]; !exists {
				g.players[state.Name] = &Player{X: 0, Y: 0, states: make([]server.State, 0)}
			}
			g.players[state.Name].AddSnapshot(state)
		}
	}

	for _, p := range g.players {
		p.NextState()
	}
	g.handleMovement()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.players {
		vector.DrawFilledRect(screen, p.X, p.Y, 5, 5, color.RGBA{255, 0, 0, 255}, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) handleMovement() {

	//Move right
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if !g.keys.wasPressing(server.MoveRight) {
			g.keys.keyPress(server.MoveRight)
			g.manager.PerformAction(server.Action{Input: server.MoveRight, Player: "John"})
		}
	} else if g.keys.wasPressing(server.MoveRight) {
		g.keys.keyRelease(server.MoveRight)
		g.manager.PerformAction(server.Action{Input: server.StopMoveRight, Player: "John"})
	}

	//Move down
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if !g.keys.wasPressing(server.MoveDown) {
			g.keys.keyPress(server.MoveDown)
			g.manager.PerformAction(server.Action{Input: server.MoveDown, Player: "John"})
		}
	} else if g.keys.wasPressing(server.MoveDown) {
		g.keys.keyRelease(server.MoveDown)
		g.manager.PerformAction(server.Action{Input: server.StopMoveDown, Player: "John"})
	}

	//Move left
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if !g.keys.wasPressing(server.MoveLeft) {
			g.keys.keyPress(server.MoveLeft)
			g.manager.PerformAction(server.Action{Input: server.MoveLeft, Player: "John"})
		}
	} else if g.keys.wasPressing(server.MoveLeft) {
		g.keys.keyRelease(server.MoveLeft)
		g.manager.PerformAction(server.Action{Input: server.StopMoveLeft, Player: "John"})
	}
	
	//Move up
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if !g.keys.wasPressing(server.MoveUp) {
			g.keys.keyPress(server.MoveUp)
			g.manager.PerformAction(server.Action{Input: server.MoveUp, Player: "John"})
		}
	} else if g.keys.wasPressing(server.MoveUp) {
		g.keys.keyRelease(server.MoveUp)
		g.manager.PerformAction(server.Action{Input: server.StopMoveUp, Player: "John"})
	}
}

func main() {
	processor := server.NewProcessor()
	manager := server.CreateNewGame(&processor)
	g := &Game{manager: &manager, players: make(map[string]*Player), keys: KeyController{keyMap: make(map[server.InputType]bool)}}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Goshooter Demo")
	quitChan := manager.Start()
	if err := ebiten.RunGame(g); err != nil {
		close(quitChan)
		log.Fatal(err)
	}
}
