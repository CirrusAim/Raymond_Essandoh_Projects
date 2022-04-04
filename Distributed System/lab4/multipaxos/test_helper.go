// Test helper functions - DO NOT EDIT

package multipaxos

type proposerTest struct {
	proposer *Proposer
	actions  []proposerAction
}

type proposerAction struct {
	promise    Promise
	desc       string
	wantOutput bool
	wantAccs   []Accept
}

type msgtype int

const (
	prepare msgtype = iota
	accept
)

type acceptorTest struct {
	acceptor *Acceptor
	actions  []acceptorAction
}

type acceptorAction struct {
	desc       string
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
	wantSid    SlotID
}

type mockLD struct{}

func (l *mockLD) Leader() int {
	return -1
}

func (l *mockLD) Subscribe() <-chan int {
	return make(chan int)
}
