package server

import (
	"github.com/ByteArena/box2d"
)

type InputType string

const (
	MoveLeft  InputType = "LEFT"
	StopMoveLeft InputType = "STOP_LEFT"
	MoveUp    InputType = "UP"
	StopMoveUp InputType = "STOP_UP"
	MoveRight InputType = "RIGHT"
	StopMoveRight InputType = "STOP_RIGHT"
	MoveDown  InputType = "DOWN"
	StopMoveDown InputType = "STOP_DOWN"
)

const moveSpeed = 100.0

type Action struct {
	Input  InputType
	Player string
}
type Processor interface {
	Process(playerStore map[string]*box2d.B2Body, actions []Action) error
	CreateCharacter() *box2d.B2Body
}
type WorldProcessor struct {
	World *box2d.B2World
}

func (e *WorldProcessor) Process(playerStore map[string]*box2d.B2Body, actions []Action) error {
	for _, action := range actions {
		player := playerStore[action.Player]
		movePlayer(action.Input, player)
	}

	tickRate := 60 //World iterations per second
	tickDuration := 1.0 / float64(tickRate)

	e.World.Step(tickDuration, 6, 2)

	return nil
}

func (e *WorldProcessor) CreateCharacter() *box2d.B2Body {
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(2.0, 2.0)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.FixedRotation = false
	bd.AllowSleep = false

	character := e.World.CreateBody(&bd)
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(0.20, 0.20)
	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = 1.0
	character.CreateFixtureFromDef(&fd)
	return character
}

func setupWorld() box2d.B2World {
	gravity := box2d.MakeB2Vec2(0, 0)
	return box2d.MakeB2World(gravity)
}

func movePlayer(input InputType, player *box2d.B2Body) {
	velocity := player.GetLinearVelocity()

	switch input {
	case MoveUp:
		player.SetLinearVelocity(box2d.B2Vec2{X: velocity.X, Y: moveSpeed})	
	case MoveLeft:
		player.SetLinearVelocity(box2d.B2Vec2{X: moveSpeed * -1, Y: velocity.Y})	
	case MoveDown:
		player.SetLinearVelocity(box2d.B2Vec2{X: velocity.X, Y: moveSpeed * -1})
	case MoveRight:
		player.SetLinearVelocity(box2d.B2Vec2{X: moveSpeed, Y: velocity.Y})
	case StopMoveRight, StopMoveLeft:
		player.SetLinearVelocity(box2d.B2Vec2{X: 0.0, Y: velocity.Y})
	case StopMoveUp, StopMoveDown:
		player.SetLinearVelocity(box2d.B2Vec2{X: velocity.X, Y: 0})
	
	}
}

func NewProcessor() WorldProcessor {
	world := setupWorld()
	return WorldProcessor{World: &world}
}

