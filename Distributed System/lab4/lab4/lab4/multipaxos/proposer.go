package multipaxos

import (
	"dat520/lab3/leaderdetector"
	"sort"
)

// Proposer represents a proposer as defined by the Multi-Paxos algorithm.
type Proposer struct {
	id           int
	quorum       int
	n            int
	crnd         Round
	adu          SlotID
	nextSlot     SlotID
	promises     []*Promise
	promiseCount int
	ld           leaderdetector.LeaderDetector
	leader       int
}

// NewProposer returns a new Multi-Paxos proposer. It takes the following
// arguments:
//
// id: The id of the node running this instance of a Paxos proposer.
//
// nrOfNodes: The total number of Paxos nodes.
//
// adu: all-decided-up-to. The initial id of the highest _consecutive_ slot
// that has been decided. Should normally be set to -1 initially, but for
// testing purposes it is passed in the constructor.
//
// ld: A leader detector implementing the detector.LeaderDetector interface.
//
// The proposer's internal crnd field should initially be set to the value of
// its id.
func NewProposer(id, nrOfNodes, adu int, ld leaderdetector.LeaderDetector) *Proposer {
	return &Proposer{
		id:       id,
		quorum:   (nrOfNodes / 2) + 1,
		n:        nrOfNodes,
		crnd:     Round(id),
		adu:      SlotID(adu),
		nextSlot: 0,
		promises: make([]*Promise, nrOfNodes),
		ld:       ld,
		leader:   ld.Leader(),
	}
}

// handlePromise processes promise prm according to the Multi-Paxos
// algorithm. If handling the promise results in proposer p emitting a
// corresponding accept slice, then output will be true and accs contain the
// accept messages. If handlePromise returns false as output, then accs will be
// a nil slice.
func (p *Proposer) handlePromise(prm Promise) (accs []Accept, output bool) {
	// TODO(student): algorithm implementation
	if prm.Rnd != p.crnd || contains(p.promises, prm.From) {
		return nil, false
	}

	p.promises = append(p.promises, &prm)
	p.nextSlot++
	if int(p.nextSlot) < p.quorum {
		return nil, false
	}

	var mp = make(map[SlotID]PromiseSlot)
	for _, prm := range p.promises {
		if prm == nil {
			continue
		}
		for _, s := range prm.Slots {
			v, ok := mp[s.ID]
			if ok {
				if s.Vrnd > v.Vrnd {
					mp[s.ID] = s
				}
			} else {
				mp[s.ID] = s
			}
		}
	}

	keys := make([]int, 0, len(mp))
	for k := range mp {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	final := make([]SlotID, 0)
	idx := 0
	if len(keys) > 0 {
		for i := keys[0]; i <= keys[len(keys)-1]; i++ {
			if keys[idx] == i {
				final = append(final, SlotID(keys[idx]))
				idx++
			} else {
				final = append(final, SlotID(i))
			}
		}
	}

	accepts := make([]Accept, 0)
	for _, k := range final {
		if p.adu >= k {
			continue
		}
		v, ok := mp[k]
		if !ok {
			accepts = append(accepts, Accept{
				From: p.id,
				Slot: k,
				Rnd:  p.crnd,
				Val:  Value{Noop: true},
			})
		} else {
			accepts = append(accepts, Accept{
				From: p.id,
				Slot: k,
				Rnd:  p.crnd,
				Val:  v.Vval,
			})
		}
	}

	return accepts, true
}

func contains(promises []*Promise, currId int) bool {
	for _, prm := range promises {
		if prm == nil {
			continue
		}
		if prm.From == currId {
			return true
		}
	}
	return false
}
