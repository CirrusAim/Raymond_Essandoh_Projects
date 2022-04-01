package singlepaxos

// Proposer represents a proposer as defined by the single-decree Paxos
// algorithm.
type Proposer struct {
	crnd        Round
	clientValue Value
	// TODO(student): algorithm implementation
	// Add other needed fields
	id              int
	totalPaxosNodes int
	acceptedNodes   map[Round][]int
	threshold       int
	maxVrnd         Round
}

// NewProposer returns a new single-decree Paxos proposer.
// It takes the following arguments:
//
// id: The id of the node running this instance of a Paxos proposer.
//
// nrOfNodes: The total number of Paxos nodes.
//
// The proposer's internal crnd field should initially be set to the value of
// its id.
func NewProposer(id int, nrOfNodes int) *Proposer {
	// TODO(student): algorithm implementation

	p := &Proposer{
		id:              id,
		crnd:            Round(id),
		totalPaxosNodes: nrOfNodes,
		acceptedNodes:   make(map[Round][]int),
		maxVrnd:         NoRound,
	}

	p.threshold = (nrOfNodes / 2) + 1

	return p
}

// handlePromise processes promise prm according to the single-decree
// Paxos algorithm. If handling the promise results in proposer p emitting a
// corresponding accept, then output will be true and acc contain the promise.
// If handlePromise returns false as output, then acc will be a zero-valued
// struct.
func (p *Proposer) handlePromise(prm Promise) (acc Accept, output bool) {
	// TODO(student): algorithm implementation

	// Check if Round is different
	if prm.Rnd != p.crnd {
		return Accept{}, false
	}

	// check if other node has aleady proposed for that Round
	if contains(p.acceptedNodes[p.crnd], prm.From) {
		return Accept{}, false
	}

	p.acceptedNodes[p.crnd] = append(p.acceptedNodes[p.crnd], prm.From)
	if p.maxVrnd < prm.Vrnd {
		p.maxVrnd = prm.Vrnd
		p.clientValue = prm.Vval
	}

	if len(p.acceptedNodes[p.crnd]) < p.threshold {
		return Accept{}, false
	} else {
		return Accept{From: p.id, Rnd: p.crnd, Val: p.clientValue}, true
	}
}

// increaseCrnd increases proposer p's crnd field by the total number
// of Paxos nodes.
func (p *Proposer) increaseCrnd() {
	// TODO(student): algorithm implementation
	p.crnd += Round(p.totalPaxosNodes)
}

// TODO(student): Add any other unexported methods needed.
func contains(source []int, val int) bool {
	for _, v := range source {
		if v == val {
			return true
		}
	}
	return false
}
