package main

type ActionQueue struct {
	actions []Action
}

func (aq *ActionQueue) Push(action Action) {
	aq.actions = append(aq.actions, action)
}

func (aq *ActionQueue) PopAll() []Action {
	actionsBuffer := make([]Action, len(aq.actions))
	copy(actionsBuffer, aq.actions)
	aq.actions = aq.actions[:0]
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