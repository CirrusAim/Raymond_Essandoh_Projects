package singlepaxos

// Learner represents a learner as defined by the single-decree Paxos
// algorithm.
type Learner struct { // TODO(student): algorithm implementation
	// Add needed fields
	// Tip: you need to keep the decided values by the Paxos nodes somewhere
	learners        map[Round][]int
	id              int
	totalPaxosNodes int
	currLearnerRnd  Round
	threshold       int
}

// NewLearner returns a new single-decree Paxos learner. It takes the
// following arguments:
//
// id: The id of the node running this instance of a Paxos learner.
//
// nrOfNodes: The total number of Paxos nodes.
func NewLearner(id int, nrOfNodes int) *Learner {
	// TODO(student): algorithm implementation
	th := (nrOfNodes / 2) + 1
	return &Learner{
		id:              id,
		totalPaxosNodes: nrOfNodes,
		learners:        make(map[Round][]int),
		currLearnerRnd:  NoRound,
		threshold:       th,
	}
}

// handleLearn processes learn lrn according to the single-decree
// Paxos algorithm. If handling the learn results in learner l emitting a
// corresponding decided value, then output will be true and val contain the
// decided value. If handleLearn returns false as output, then val will have
// its zero value.
func (l *Learner) handleLearn(learn Learn) (val Value, output bool) {
	// TODO(student): algorithm implementation

	// if incoming Round is less than the max round, ignore
	if learn.Rnd < l.currLearnerRnd {
		return "", false
	} else if learn.Rnd == l.currLearnerRnd {
		if !contains(l.learners[learn.Rnd], learn.From) {
			l.learners[learn.Rnd] = append(l.learners[learn.Rnd], learn.From)
		}
		l.currLearnerRnd = learn.Rnd
	} else {
		l.learners[learn.Rnd] = append(l.learners[learn.Rnd], learn.From)
		l.currLearnerRnd = learn.Rnd
	}

	if len(l.learners[l.currLearnerRnd]) < l.threshold {
		return "", false
	}

	return learn.Val, true
}

// TODO(student): Add any other unexported methods needed.
