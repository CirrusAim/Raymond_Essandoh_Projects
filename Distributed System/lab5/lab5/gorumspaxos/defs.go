package gorumspaxos

import (
	"dat520/lab3/leaderdetector"
	pb "dat520/lab5/gorumspaxos/proto"
	"time"
)

func NoRound() *pb.Round {
	return &pb.Round{Id: -1}
}

type msgtype int

// Constants used in the acceptor tests
const (
	prepare msgtype = iota
	accept
)

// NewProposerArgs holds the arguments to initialize the Proposer
type NewProposerArgs struct {
	id               int                           // id: The id of the node running this instance of a Paxos proposer.
	aduSlotID        uint32                        // aduSlotId: all-decided-up-to. The initial id of the highest _consecutive_ slot
	leaderDetector   leaderdetector.LeaderDetector // leaderDetector: A leader detector implementing the detector.LeaderDetector interface.
	qspec            pb.QuorumSpec                 // qspec: Implementor of the quorum spec.
	nodeMap          map[string]uint32             // nodeMap contains the map of address to id
	phaseOneWaitTime time.Duration                 // phaseOneWaitTime: duration until the timeout for phase1 operations
	phaseTwoWaitTime time.Duration                 // phaseTwoWaitTime: duration until the timeout for phase2 operations
}

type NewPaxosReplicaArgs struct {
	LocalAddr string            // LocalAddr: is the address of the replica
	Id        int               // Id is the id of the current node
	NodeMap   map[string]uint32 // NodeIds: contains the the node id and the addresses of the replicas
}

// struct acceptorAction is used in acceptor testing
type acceptorAction struct {
	desc    string
	msgtype msgtype
	prepare *pb.PrepareMsg
	accept  *pb.AcceptMsg
	wantPrm *pb.PromiseMsg
	wantLrn *pb.LearnMsg
}

// comparePromiseMsgs compares the promise messages and return true if both messages
// contains the same values.
func comparePromiseMsgs(a, b *pb.PromiseMsg) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil {
		return false
	} else if b == nil {
		return false
	}
	return checkPromiseSlots(a.Slots, b.Slots) && checkIsRoundEqual(a.Rnd, b.Rnd)
}

// checkPromiseSlots compares the promise slots and return true if both structs
// contains the same values.
func checkPromiseSlots(a, b []*pb.PromiseSlot) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil && b != nil {
		return false
	} else if a != nil && b == nil {
		return false
	} else if len(a) != len(b) {
		return false
	}
	for _, slot := range a {
		if !isSlotInSlice(slot, b) {
			return false
		}
	}
	return true
}

// isSlotInSlice: checks if a given slot is present in the slice of promise slots.
func isSlotInSlice(a *pb.PromiseSlot, b []*pb.PromiseSlot) bool {
	for _, slot := range b {
		if checkIsPromiseSlotEqual(slot, a) {
			return true
		}
	}
	return false
}

// checkIsPromiseSlotEqual: checks if two promise slot structures are same by
// comparing the values of the elements.
func checkIsPromiseSlotEqual(a, b *pb.PromiseSlot) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil {
		return false
	} else if b == nil {
		return false
	} else if !checkIsSlotEqual(a.Slot, b.Slot) {
		return false
	} else if !checkIsRoundEqual(a.Vrnd, b.Vrnd) {
		return false
	} else if !checkIsValueEqual(a.Value, b.Value) {
		return false
	}
	return true
}

// checkIsSlotEqual: checks if two slot structures are same by
// comparing the values of the elements.
func checkIsSlotEqual(a, b *pb.Slot) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil {
		return false
	} else if b == nil {
		return false
	} else if a.Id != b.Id {
		return false
	}
	return true
}

// checkIsRoundEqual: checks if two round structures are same by
// comparing the values of the elements
func checkIsRoundEqual(a, b *pb.Round) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil {
		return false
	} else if b == nil {
		return false
	} else if a.Id != b.Id {
		return false
	}
	return true
}

// checkIsValueEqual: checks if two value structures are same by
// comparing the values of the elements
func checkIsValueEqual(a, b *pb.Value) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil {
		return false
	} else if b == nil {
		return false
	} else if a.ClientCommand != b.ClientCommand {
		return false
	} else if a.ClientID != b.ClientID {
		return false
	} else if a.IsNoop != b.IsNoop {
		return false
	} else if a.ClientSeq != b.ClientSeq {
		return false
	}
	return true
}

// struct proposeQFTest is useful in testing the proposer implementation.
type proposeQFTest struct {
	quorumSize       int
	promiseMsgs      []*pb.PromiseMsg
	inputPrepareMsg  *pb.PrepareMsg
	outputPromiseMsg *pb.PromiseMsg
	wantOutput       bool
	description      string
}

// struct acceptQFTest is useful in testing the acceptor implementation.
type acceptQFTest struct {
	description    string
	quorumSize     int
	learnMsgs      []*pb.LearnMsg
	inputAcceptMsg *pb.AcceptMsg
	outputLearnMsg *pb.LearnMsg
	wantOutput     bool
}

type clientHandleQFTest struct {
	description string
	inputValue  *pb.Value
	QuorumSize  int
	replies     []*pb.Response
	wantOutput  bool
}

type commitQFTest struct {
	description    string
	quorumSize     int
	testQuorumSize int
	wantOutput     bool
}

// Example values to be used in testing of acceptor and proposer.
var (
	testingValueOne = pb.Value{
		ClientID:      "1234",
		ClientSeq:     42,
		ClientCommand: "ls",
	}
	testingValueTwo = pb.Value{
		ClientID:      "5678",
		ClientSeq:     99,
		ClientCommand: "rm",
	}
	testingValueThree = pb.Value{
		ClientID:      "1369",
		ClientSeq:     4,
		ClientCommand: "mkdir",
	}
)
