package leaderdetector

// A MonLeaderDetector represents a Monarchical Eventual Leader Detector as
// described at page 53 in:
// Christian Cachin, Rachid Guerraoui, and Luís Rodrigues: "Introduction to
// Reliable and Secure Distributed Programming" Springer, 2nd edition, 2011.
type MonLeaderDetector struct {
	nodeIDs     []int
	leader      int
	suspected   map[int]bool
	subscribers []chan int
}

// NewMonLeaderDetector returns a new Monarchical Eventual Leader Detector
// given a list of node ids.
func NewMonLeaderDetector(nodeIDs []int) *MonLeaderDetector {
	suspected := make(map[int]bool)
	subscriber := make([]chan int, 0)

	m := &MonLeaderDetector{
		nodeIDs:     nodeIDs,
		suspected:   suspected,
		subscribers: subscriber,
	}
	m.leader = m.Leader()
	return m
}

// Leader returns the current leader. Leader will return UnknownID if all nodes
// are suspected.
func (m *MonLeaderDetector) Leader() int {
	maxId := UnknownID
	for _, node := range m.nodeIDs {
		if _, ok := m.suspected[node]; ok {
			continue
		}
		if node > maxId {
			maxId = node
		}
	}
	return maxId
}

// Suspect instructs the leader detector to consider the node with matching
// id as suspected. If the suspect indication result in a leader change
// the leader detector should publish this change to its subscribers.
func (m *MonLeaderDetector) Suspect(id int) {
	currLeader := m.Leader()
	m.suspected[id] = true
	newLeader := m.Leader()
	if newLeader != currLeader {
		m.subscribe(newLeader)
	}
}

// Restore instructs the leader detector to consider the node with matching
// id as restored. If the restore indication result in a leader change
// the leader detector should publish this change to its subscribers.
func (m *MonLeaderDetector) Restore(id int) {
	currLeader := m.Leader()
	delete(m.suspected, id)
	newLeader := m.Leader()
	if newLeader != currLeader {
		m.subscribe(newLeader)
	}
}

// Subscribe returns a buffered channel which will be used by the leader
// detector to publish the id of the highest ranking node.
// The leader detector will publish UnknownID if all nodes become suspected.
// Subscribe will drop publications to slow subscribers.
// Note: Subscribe returns a unique channel to every subscriber;
// it is not meant to be shared.
func (m *MonLeaderDetector) Subscribe() <-chan int {
	subscriber := make(chan int, 1)
	m.subscribers = append(m.subscribers, subscriber)
	return subscriber
}

func (m *MonLeaderDetector) subscribe(leader int) {
	for _, v := range m.subscribers {
		v <- leader
	}
}
