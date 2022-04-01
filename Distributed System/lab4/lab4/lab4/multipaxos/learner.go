package multipaxos

// Learner represents a learner as defined by the Multi-Paxos algorithm.
type Learner struct {
	// TODO(student): algorithm implementation
	// Add needed fields
	learners        map[SlotID]map[int]Learn
	id              int
	totalPaxosNodes int
	currLearnerRnd  Round
	threshold       int
}

// NewLearner returns a new Multi-Paxos learner. It takes the
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
		learners:        make(map[SlotID]map[int]Learn),
		currLearnerRnd:  NoRound,
		threshold:       th,
	}
}

// handleLearn processes learn lrn according to the Multi-Paxos
// algorithm. If handling the learn results in learner l emitting a
// corresponding decided value, then output will be true, sid the id for the
// slot that was decided and val contain the decided value. If handleLearn
// returns false as output, then val and sid will have their zero value.
func (l *Learner) handleLearn(learn Learn) (val Value, sid SlotID, output bool) {
	// TODO(student): algorithm implementation
	v, ok := l.learners[learn.Slot]
	if ok {
		frm, fok := v[learn.From]
		if fok {
			if frm.Rnd <= learn.Rnd {
				return Value{}, 0, false
			} else {
				v[learn.From] = learn
			}
		} else {
			l.learners[learn.Slot][learn.From] = learn
		}

	} else {
		l.learners[learn.Slot] = map[int]Learn{learn.From: learn}
	}

	if len(l.learners[learn.Slot]) < l.threshold {
		return Value{}, 0, false
	} else {
		frms := l.learners[learn.Slot]
		rnd := frms[learn.From]
		qrmCnt := 0
		for _, v := range frms {
			if v.Rnd == rnd.Rnd {
				qrmCnt++
			}
		}
		if qrmCnt < l.threshold {
			return Value{}, 0, false
		}
	}
	delete(l.learners, learn.Slot)
	return learn.Val, learn.Slot, true
}

// TODO(student): Add any other unexported methods needed.
