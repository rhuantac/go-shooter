package server

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/ByteArena/box2d"
)

type InputType string

const (
	MoveLeft      InputType = "LEFT"
	StopMoveLeft  InputType = "STOP_LEFT"
	MoveUp        InputType = "UP"
	StopMoveUp    InputType = "STOP_UP"
	MoveRight     InputType = "RIGHT"
	StopMoveRight InputType = "STOP_RIGHT"
	MoveDown      InputType = "DOWN"
	StopMoveDown  InputType = "STOP_DOWN"
	Rotate        InputType = "ROTATE"
	Shoot         InputType = "SHOOT"
)

const moveSpeed = 10.0

type Action struct {
	Input  InputType
	Angle  float64 `json:",omitempty"`
	Player string
}

type BodyData struct {
	ID       string
	Type     string
	UserData interface{}
}

type Processor interface {
	Process(playerStore map[string]*box2d.B2Body, objectStore map[string]*box2d.B2Body, actions []Action) error
	CreateCharacter() *box2d.B2Body
}
type WorldProcessor struct {
	World *box2d.B2World
}

func (e *WorldProcessor) Process(playerStore map[string]*box2d.B2Body, objectStore map[string]*box2d.B2Body, actions []Action) error {
	for _, action := range actions {
		fmt.Printf("action %v", action)
		player := playerStore[action.Player]
		if player == nil {
			return fmt.Errorf("player %s not found", action.Player)
		}
		switch action.Input {
		case Rotate:
			player.SetTransform(player.GetPosition(), action.Angle)
		case MoveLeft, MoveRight, MoveUp, MoveDown, StopMoveLeft, StopMoveRight, StopMoveUp, StopMoveDown:
			movePlayer(action.Input, player)
		case Shoot:
			bullet := shoot(player)
			bulletID := fmt.Sprintf("bullet-%d", rand.Int())
			bulletData := BodyData{ID: bulletID, Type: "bullet", UserData: map[string]interface{}{"owner": action.Player}}
			bullet.SetUserData(bulletData)
			objectStore[bulletID] = bullet
		}
	}

	tickRate := 30 //World iterations per second
	tickDuration := 1.0 / float64(tickRate)

	e.World.Step(tickDuration, 6, 2)
	contactEvents := e.World.GetContactList()
	for contactEvents != nil {
		bodyA := contactEvents.GetFixtureA().GetBody()
		bodyB := contactEvents.GetFixtureB().GetBody()
		userDataA := bodyA.GetUserData()
		userDataB := bodyB.GetUserData()
		if userDataA != nil && userDataB != nil {
			dataA := userDataA.(BodyData)
			dataB := userDataB.(BodyData)
			if dataA.Type == "bullet" || dataB.Type == "bullet" {
				// Remove bullet from world and objectStore
				if dataA.Type == "bullet" {
					//Check if body exists before destroying
					e.World.DestroyBody(bodyA)
					delete(objectStore, dataA.ID)
				}
				if dataB.Type == "bullet" {
					e.World.DestroyBody(bodyB)
					delete(objectStore, dataB.ID)
				}
			}
		}
		contactEvents = contactEvents.GetNext()
	}
	return nil
}

func shoot(player *box2d.B2Body) *box2d.B2Body {
	const armLength float64 = 9.0
	const shotSpeed = 50
	bd := box2d.MakeB2BodyDef()
	spawnY, spawnX := math.Sincos(player.GetAngle())
	bd.Position.Set(player.GetPosition().X+spawnX*armLength, player.GetPosition().Y+spawnY*armLength)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.FixedRotation = true
	bd.AllowSleep = false
	bd.Bullet = true

	bullet := player.GetWorld().CreateBody(&bd)
	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(1)
	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = 1.0
	bullet.CreateFixtureFromDef(&fd)
	velocityX := shotSpeed * spawnX

	velocityY := shotSpeed * spawnY

	bullet.SetLinearVelocity(box2d.B2Vec2{X: float64(velocityX), Y: float64(velocityY)})
	return bullet
}

func (e *WorldProcessor) CreateCharacter() *box2d.B2Body {
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(20.0, 20.0)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.FixedRotation = true
	bd.AllowSleep = false

	character := e.World.CreateBody(&bd)
	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(5)
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
		player.SetLinearVelocity(box2d.B2Vec2{X: velocity.X, Y: moveSpeed * -1})
	case MoveLeft:
		player.SetLinearVelocity(box2d.B2Vec2{X: moveSpeed * -1, Y: velocity.Y})
	case MoveDown:
		player.SetLinearVelocity(box2d.B2Vec2{X: velocity.X, Y: moveSpeed})
	case MoveRight:
		player.SetLinearVelocity(box2d.B2Vec2{X: moveSpeed, Y: velocity.Y})
	case StopMoveRight:
		if velocity.X > 0 { // Only stop right movement
			player.SetLinearVelocity(box2d.B2Vec2{X: 0, Y: velocity.Y})
		}
	case StopMoveLeft:
		if velocity.X < 0 { // Only stop left movement
			player.SetLinearVelocity(box2d.B2Vec2{X: 0, Y: velocity.Y})
		}
	case StopMoveUp:
		if velocity.Y < 0 { // Only stop up movement
			player.SetLinearVelocity(box2d.B2Vec2{X: velocity.X, Y: 0})
		}
	case StopMoveDown:
		if velocity.Y > 0 { // Only stop down movement
			player.SetLinearVelocity(box2d.B2Vec2{X: velocity.X, Y: 0})
		}
	}
}

func NewProcessor() WorldProcessor {
	world := setupWorld()
	return WorldProcessor{World: &world}
}
