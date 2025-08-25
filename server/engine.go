package server

import (
	"fmt"
	"log"
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
	Hit           InputType = "HIT"
	Destroy       InputType = "DESTROY"
)
const moveSpeed = 10.0

type Action struct {
	Input      InputType
	Angle      float64 `json:",omitempty"`
	Player     string
	FromServer bool   `json:"fromServer,omitempty"`
	Target     string `json:"target,omitempty"`
}

type ServerEvent struct {
	Type   string      `json:"type"`
	Body   interface{} `json:"body"`
	Player string      `json:"player,omitempty"`
}

type BodyData struct {
	ID       string
	Type     string
	UserData interface{}
}

type Processor interface {
	Process(playerStore map[string]*box2d.B2Body, objectStore map[string]*box2d.B2Body, actions []Action) ([]ServerEvent, error)
	CreateCharacter() *box2d.B2Body
}
type WorldProcessor struct {
	World *box2d.B2World
}

func (e *WorldProcessor) Process(playerStore map[string]*box2d.B2Body, objectStore map[string]*box2d.B2Body, actions []Action) ([]ServerEvent, error) {
	serverActions := []ServerEvent{}
	for _, action := range actions {
		fmt.Printf("action %v", action)
		player := playerStore[action.Player]
		if player == nil && !action.FromServer {
			return nil, fmt.Errorf("player %s not found", action.Player)
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
		case Destroy:
			target := objectStore[action.Target]
			if target != nil {
				log.Printf("Destroying object %s", action.Target)
				e.World.DestroyBody(objectStore[action.Target])
				delete(objectStore, action.Target)
			}
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
					serverActions = append(serverActions, ServerEvent{Type: "HIT", Body: map[string]string{"bulletId": dataA.ID, "target": dataB.ID, "posX": fmt.Sprintf("%f", bodyA.GetPosition().X), "posY": fmt.Sprintf("%f", bodyA.GetPosition().Y)}})
				}
				if dataB.Type == "bullet" {
					serverActions = append(serverActions, ServerEvent{Type: "HIT", Body: map[string]string{"bulletId": dataB.ID, "target": dataA.ID, "posX": fmt.Sprintf("%f", bodyB.GetPosition().X), "posY": fmt.Sprintf("%f", bodyB.GetPosition().Y)}})
				}
			}
		}
		contactEvents = contactEvents.GetNext()
	}
	return serverActions, nil
}

func shoot(player *box2d.B2Body) *box2d.B2Body {
	const armLength float64 = 9.0
	const shotSpeed = 20
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
	bd.Position.Set(rand.Float64()*50, rand.Float64()*50)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.FixedRotation = true
	bd.AllowSleep = true

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
