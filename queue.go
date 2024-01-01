package main

import "sync"

type ActionQueue struct {
	actions []Action
	mutex *sync.Mutex
}

func (aq *ActionQueue) Size () int{
	aq.mutex.Lock()
	defer aq.mutex.Unlock()
	return len(aq.actions)
}
func (aq *ActionQueue) Push(action Action) {
	aq.mutex.Lock()
	aq.actions = append(aq.actions, action)
	aq.mutex.Unlock()
}

func (aq *ActionQueue) PopAll() []Action {
	aq.mutex.Lock()
	actionsBuffer := make([]Action, len(aq.actions))
	copy(actionsBuffer, aq.actions)
	aq.actions = aq.actions[:0]
	aq.mutex.Unlock()
	return actionsBuffer
}

type SnapshotQueue struct {
	snapshots []Snapshot
}

func (q *SnapshotQueue) Push(snapshot Snapshot) {
	q.snapshots = append(q.snapshots, snapshot)
}

func (q *SnapshotQueue) PopAll() []Snapshot {
	snapshotBuffer := make([]Snapshot, len(q.snapshots))
	copy(snapshotBuffer, q.snapshots)
	q.snapshots = q.snapshots[:0]
	return snapshotBuffer
}

func NewActionQueue() *ActionQueue{
	return &ActionQueue{actions: []Action{}, mutex: &sync.Mutex{}}
}