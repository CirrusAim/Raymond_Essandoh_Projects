package singlepaxos

// Acceptor represents an acceptor as defined by the single-decree Paxos
// algorithm.
type Acceptor struct { // TODO(student): algorithm implementation
	// Add needed fields
	prvPrepareRnd    Round
	prvAcceptVal     Value
	isAnyValAccepted bool
	prvAcceptVrnd    Round
	id               int
}

// NewAcceptor returns a new single-decree Paxos acceptor.
// It takes the following arguments:
//
// id: The id of the node running this instance of a Paxos acceptor.
func NewAcceptor(id int) *Acceptor {
	// TODO(student): algorithm implementation
	return &Acceptor{
		id:               id,
		prvAcceptVal:     ZeroValue,
		isAnyValAccepted: false,
		prvPrepareRnd:    NoRound,
		prvAcceptVrnd:    NoRound,
	}
}

// handlePrepare processes prepare prp according to the single-decree
// Paxos algorithm. If handling the prepare results in acceptor a emitting a
// corresponding promise, then output will be true and prm contain the promise.
// If handlePrepare returns false as output, then prm will be a zero-valued
// struct.
func (a *Acceptor) handlePrepare(prp Prepare) (prm Promise, output bool) {
	// TODO(student): algorithm implementation
	if prp.Crnd < a.prvPrepareRnd {
		return Promise{}, false
	}
	a.prvPrepareRnd = prp.Crnd

	if a.isAnyValAccepted {
		p := Promise{
			To:   prp.From,
			From: a.id,
			Rnd:  prp.Crnd,
			Vrnd: a.prvAcceptVrnd,
			Vval: a.prvAcceptVal,
		}
		return p, true
	}

	return Promise{To: prp.From, From: a.id, Rnd: prp.Crnd, Vrnd: NoRound, Vval: ZeroValue}, true
}

// handleAccept processes accept acc according to the single-decree
// Paxos algorithm. If handling the accept results in acceptor a emitting a
// corresponding learn, then output will be true and lrn contain the learn.  If
// handleAccept returns false as output, then lrn will be a zero-valued struct.
func (a *Acceptor) handleAccept(acc Accept) (lrn Learn, output bool) {
	// TODO(student): algorithm implementation
	if acc.Rnd < a.prvPrepareRnd {
		return Learn{}, false
	}

	a.isAnyValAccepted = true
	a.prvAcceptVal = acc.Val
	a.prvAcceptVrnd = acc.Rnd
	return Learn{From: a.id, Rnd: acc.Rnd, Val: acc.Val}, true
}

// TODO(student): Add any other unexported methods needed.
