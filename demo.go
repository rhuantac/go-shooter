package main

import (
	"errors"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	screenWidth  = 480
	screenHeight = 480
)

type Player struct {
	X, Y   float32
	states []State
}

func (p *Player) AddSnapshot(s State) {
	if len(p.states) >= 60 {
		// keep only the last 10 snaps
		p.states = append(p.states[:60], s)
		return
	}

	p.states = append(p.states, s)
}

func (p *Player) RemoveSnapshot() (State, error) {
	if len(p.states) == 0 {
		return State{}, errors.New("array is empty")
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
	p.X = float32(next.posX)
	p.Y = float32(-1 * next.posY) //Server coordinates are cartesian
}

type Game struct {
	manager *GameManager
	players map[string]*Player
}

func (g *Game) Update() error {
	snaps := g.manager.GetSnapshots()
	for _, snap := range snaps {
		for _, state := range snap.state {
			if _, exists := g.players[state.name]; !exists {
				g.players[state.name] = &Player{X: 0, Y: 0, states: make([]State, 0)}
			}
			g.players[state.name].AddSnapshot(state)
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
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.manager.PerformAction(Action{input: MoveRight, player: "John"})
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.manager.PerformAction(Action{input: MoveDown, player: "John"})
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.manager.PerformAction(Action{input: MoveLeft, player: "John"})
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.manager.PerformAction(Action{input: MoveUp, player: "John"})
	}
}

func main() {
	processor := NewProcessor()
	manager := CreateNewGame(&processor)
	g := &Game{manager: &manager, players: make(map[string]*Player)}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Goshooter Demo")
	quitChan := manager.Start()
	if err := ebiten.RunGame(g); err != nil {
		close(quitChan)
		log.Fatal(err)
	}
}
