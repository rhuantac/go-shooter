package main

import (
	"testing"
	"time"
)

func TestOperations(t *testing.T) {
	t.Run("actions are pushed to queue correctly", func(t *testing.T) {
		queue := NewActionQueue()
		for i := 0; i < 3; i ++ {
			queue.Push(Action{input: MoveDown, player: "John"})
		}

		if len(queue.actions) != 3 {
			t.Errorf("got %d queue actions want %d", len(queue.actions), 3)
		}

	})

	t.Run("queue is empty when actions are popped", func(t *testing.T) {
		queue := NewActionQueue()
		for i := 0; i < 3; i ++ {
			queue.Push(Action{input: MoveDown, player: "John"})
		}
		queue.PopAll()
		
		if len(queue.actions) != 0 {
			t.Errorf("got %d queue actions want %d", len(queue.actions), 3)
		}

	})

	t.Run("queue can be used concurrently", func(t *testing.T) {
		queue := NewActionQueue()
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