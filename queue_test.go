package main

import (
	"testing"
	"time"
)

func TestQueues(t *testing.T) {
	t.Run("action queue", func(t *testing.T) {
		queue := NewQueue[Action]()
		for i := 0; i < 3; i ++ {
			queue.Push(Action{input: MoveDown, player: "John"})
		}

		if queue.Size() != 3 {
			t.Errorf("got %d queue actions want %d", queue.Size(), 3)
		}

		queue.PopAll()
		
		if queue.Size() != 0 {
			t.Errorf("got %d queue actions want %d", queue.Size(), 3)
		}

		//test concurrency
		go func(){
			for i := 0; i < 10; i++ {
				queue.Push(Action{input: MoveDown, player: "John"})
			}
		}()
		go func(){
			for i := 0; i < 10; i++ {
				queue.Push(Action{})
			}
		}()

		time.Sleep(10 * time.Millisecond)
		if queue.Size() !=20 {
			t.Errorf("got %d queue actions want %d", queue.Size(), 20)
		}
	})
}