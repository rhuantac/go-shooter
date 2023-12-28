package main

import "github.com/ByteArena/box2d"

type PlayerStore map[string]*box2d.B2Body
type GameManager struct {
	gameProcessor Processor
	playerStore PlayerStore
}

func CreateNewGame(processor Processor) GameManager{
	store := make(PlayerStore, 0)
	store["John"] = processor.CreateCharacter()
	return GameManager{gameProcessor: processor, playerStore: store}
}

func gameLoop(processor Processor, players PlayerStore) {
	for i := 0; i < 60; i++ {
		processor.Process(players, Action{input: MoveDown, player: "John"})
	}
}