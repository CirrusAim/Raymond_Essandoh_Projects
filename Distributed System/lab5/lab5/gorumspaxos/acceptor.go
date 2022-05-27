package gorumspaxos

import (
	pb "dat520/lab5/gorumspaxos/proto"
	"sort"
)

// Acceptor represents an acceptor as defined by the Multi-Paxos algorithm.
// For current implementation alpha is assumed to be 1. In future if the alpha is
// increased, then there would be significant changes to the current implementation.
type Acceptor struct {
	rnd           *pb.Round                  // rnd: is the current round in which the acceptor is participating
	slots         map[uint32]*pb.PromiseSlot // slots: is the internal data structure maintained by the acceptor to remember the slots
	maxSeenSlotId uint32                     // maxSeenSlotId: is the highest slot for which the prepare is received
}

// NewAcceptor returns a new Multi-Paxos acceptor
func NewAcceptor() *Acceptor {
	return &Acceptor{
		rnd:   NoRound(),
		slots: make(map[uint32]*pb.PromiseSlot),
	}
}

// handlePrepare takes a prepare message and returns a promise message according to
// the Multi-Paxos algorithm. If the prepare message is invalid, nil should be returned.
func (a *Acceptor) handlePrepare(prp *pb.PrepareMsg) (prm *pb.PromiseMsg) {
	if prp.Crnd.Id < a.rnd.Id {
		return nil
	}

	a.rnd = prp.Crnd

	if prp.Slot != nil && prp.Slot.Id > a.maxSeenSlotId {
		a.maxSeenSlotId = prp.Slot.Id
	}

	psMap := make(map[uint32]*pb.PromiseSlot)
	slotArr := make([]int, 0)
	ps := make([]*pb.PromiseSlot, 0)
	// Let's say slots will have ids 1,2,3,5,6 . But maxSeenSlotId in prepare till now is 3.
	// so the ps array in 33 line will contain slots with ids 3,5,6.
	for k, v := range a.slots {
		if k >= a.maxSeenSlotId {
			psMap[v.Slot.Id] = v
			//ps = append(ps, v)
			slotArr = append(slotArr, int(v.Slot.Id))
		}
	}

	sort.Ints(slotArr)

	for _, v := range slotArr {
		ps = append(ps, psMap[uint32(v)])
	}

	return &pb.PromiseMsg{Rnd: prp.Crnd, Slots: ps}
}

// handleAccept takes an accept message and returns a learn message according to
// the Multi-Paxos algorithm. If the accept message is invalid, nil should be returned.
func (a *Acceptor) handleAccept(acc *pb.AcceptMsg) (lrn *pb.LearnMsg) {
	if acc.Rnd.Id < a.rnd.Id {
		return nil
	}
	a.rnd = acc.Rnd

	a.slots[acc.Slot.Id] = &pb.PromiseSlot{
		Slot:  acc.Slot,
		Vrnd:  acc.Rnd,
		Value: acc.Val,
	}
	return &pb.LearnMsg{Rnd: acc.Rnd, Slot: acc.Slot, Val: acc.Val}
}
