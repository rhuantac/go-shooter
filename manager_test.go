package main

import (
	"testing"

	"github.com/ByteArena/box2d"
)

type stubWorldProcessor struct {
	createCharacterCalls int
	processCalls         int
}

func (w *stubWorldProcessor) Process(playerStore map[string]*box2d.B2Body, action []Action) error {
	w.processCalls += 1
	return nil
}

func (w *stubWorldProcessor) CreateCharacter() *box2d.B2Body {
	w.createCharacterCalls += 1
	return nil
}

func TestGameManager(t *testing.T) {
	t.Run("one character is created on new game", func(t *testing.T) {
		processor := stubWorldProcessor{}
		CreateNewGame(&processor)
		if processor.createCharacterCalls != 1 {
			t.Errorf("got %d calls on CreateCharacter want %d", processor.createCharacterCalls, 1)
		}
	})

	t.Run("world is processed 60 ticks per call", func(t *testing.T) {
		processor := stubWorldProcessor{}
		playerStore := make(PlayerStore)
		queueProcess(&processor, playerStore, &ActionQueue{}, &SnapshotQueue{})
		if processor.processCalls != 60 {
			t.Errorf("got %d calls on gameLoop want at least %d", processor.processCalls, 60)
		}
	})
	
	t.Run("snapshots are taken on queue process", func(t *testing.T) {
		processor := stubWorldProcessor{}
		playerStore := make(PlayerStore)
		snapshotQueue :=  SnapshotQueue{}
		queueProcess(&processor, playerStore, &ActionQueue{}, &snapshotQueue)
		if len(snapshotQueue.snapshots) != 60 {
			t.Errorf("got %d snapshots want %d", len(snapshotQueue.snapshots), 60)
		}
	})
}
