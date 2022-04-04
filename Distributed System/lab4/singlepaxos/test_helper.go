// Test helper functions - DO NOT EDIT

package singlepaxos

const (
	valueFromClientOne = "A client command"
	valueFromClientTwo = "Another client command"
)

type proposerTest struct {
	proposer    *Proposer
	desc        string
	clientValue Value
	actions     []proposerAction
}

type proposerAction struct {
	promise    Promise
	wantOutput bool
	wantAcc    Accept
}

type msgtype int

const (
	prepare msgtype = iota
	accept
)

type acceptorTest struct {
	acceptor *Acceptor
	desc string
	actions  []acceptorAction
}

type acceptorAction struct {
	msgtype    msgtype
	prepare    Prepare
	accept     Accept
	wantOutput bool
	wantPrm    Promise
	wantLrn    Learn
}

type learnerTest struct {
	learner *Learner
	desc    string
	actions []learnerAction
}

type learnerAction struct {
	learn      Learn
	wantOutput bool
	wantVal    Value
}
