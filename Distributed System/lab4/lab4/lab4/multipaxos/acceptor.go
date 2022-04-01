package multipaxos

// Acceptor represents an acceptor as defined by the Multi-Paxos algorithm.
type Acceptor struct {
	// TODO(student): algorithm implementation
	// Add needed fields
	prvPrepareRnd    Round
	prvAcceptVal     Value
	isAnyValAccepted bool
	prvAcceptVrnd    Round
	id               int
	ps               []PromiseSlot
}

// NewAcceptor returns a new Multi-Paxos acceptor.
// It takes the following arguments:
//
// id: The id of the node running this instance of a Paxos acceptor.
func NewAcceptor(id int) *Acceptor {
	// TODO(student): algorithm implementation
	return &Acceptor{
		id:               id,
		isAnyValAccepted: false,
		prvPrepareRnd:    NoRound,
		prvAcceptVrnd:    NoRound,
		prvAcceptVal:     Value{},
		ps:               make([]PromiseSlot, 0),
	}
}

// handlePrepare processes prepare prp according to the Multi-Paxos
// algorithm. If handling the prepare results in acceptor a emitting a
// corresponding promise, then output will be true and prm contain the promise.
// If handlePrepare returns false as output, then prm will be a zero-valued
// struct.
func (a *Acceptor) handlePrepare(prp Prepare) (prm Promise, output bool) {
	// TODO(student): algorithm implementation
	if prp.Crnd < a.prvPrepareRnd {
		return Promise{}, false
	}
	a.prvPrepareRnd = prp.Crnd
	if len(a.ps) == 0 {
		return Promise{To: prp.From, From: a.id, Rnd: prp.Crnd}, true
	}
	ps := make([]PromiseSlot, 0)

	for _, v := range a.ps {
		if v.ID >= prp.Slot {
			ps = append(ps, v)
		}
	}

	return Promise{To: prp.From, From: a.id, Rnd: prp.Crnd, Slots: ps}, true
}

// handleAccept processes accept acc according to the Multi-Paxos
// algorithm. If handling the accept results in acceptor a emitting a
// corresponding learn, then output will be true and lrn contain the learn.  If
// handleAccept returns false as output, then lrn will be a zero-valued struct.
func (a *Acceptor) handleAccept(acc Accept) (lrn Learn, output bool) {
	// TODO(student): algorithm implementation
	if acc.Rnd < a.prvPrepareRnd {
		return Learn{}, false
	}
	a.prvPrepareRnd = acc.Rnd
	a.isAnyValAccepted = true
	a.prvAcceptVal = acc.Val
	a.prvAcceptVrnd = acc.Rnd
	ps := PromiseSlot{
		ID:   acc.Slot,
		Vrnd: acc.Rnd,
		Vval: acc.Val,
	}

	a.addPromiseSlot(ps)

	return Learn{From: a.id, Slot: acc.Slot, Rnd: acc.Rnd, Val: acc.Val}, true
}

// TODO(student): Add any other unexported methods needed.

func (a *Acceptor) addPromiseSlot(pr PromiseSlot) {
	for idx, v := range a.ps {
		if v.ID == pr.ID && v.Vrnd <= pr.Vrnd {
			a.ps = append(a.ps[:idx], a.ps[idx+1:]...)
			break
		}
	}
	a.ps = append(a.ps, pr)
	return
}
