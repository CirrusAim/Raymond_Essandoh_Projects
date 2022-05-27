package gorumspaxos

import (
	pb "dat520/lab6/gorumspaxos/proto"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestMPHandlePrepareAndAccept(t *testing.T) {
	for _, test := range acceptorTests {
		for _, action := range test.actions {
			switch action.msgtype {
			case prepare:
				gotPrm := test.acceptor.handlePrepare(action.prepare)
				if diff := cmp.Diff(action.wantPrm, gotPrm, protocmp.Transform()); diff != "" {
					t.Logf("Description: %v", action.desc)
					t.Errorf("handlePrepare() mismatch (-wantPrm +gotPrm):\n%s", diff)
				}
			case accept:
				gotLrn := test.acceptor.handleAccept(action.accept)
				if diff := cmp.Diff(action.wantLrn, gotLrn, protocmp.Transform()); diff != "" {
					t.Logf("Description: %v", action.desc)
					t.Errorf("handleAccept() mismatch (-wantLrn +gotLrn):\n%s", diff)
				}
			default:
				t.Error("assertion failed: unknown message type for acceptor")
			}
		}
	}
}

var acceptorTests = []struct {
	acceptor *Acceptor
	actions  []acceptorAction
}{
	{
		NewAcceptor(),
		[]acceptorAction{
			{
				desc:    "prepare slot 1, crnd 2 -> expected corresponding promise",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 2},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 2},
				},
			},
			{
				desc:    "prepare slot 1, crnd 1 -> expected nil promise, ignore due to lower crnd",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 1},
				},
				wantPrm: nil,
			},
		},
	},
	{
		NewAcceptor(),
		[]acceptorAction{
			{
				desc:    "prepare with crnd 0 and slot 1, initial acceptor round is NoRound (-1) -> expected correct promise",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 0},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 0},
				},
			},
		},
	},
	{
		NewAcceptor(),
		[]acceptorAction{
			{
				desc:    "accept for slot 1 with rnd 2, current acceptor rnd should be NoRound (-1) -> expected learn with correct slot, rnd and value",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueThree,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueThree,
				},
			},
			{
				desc:    "prepare, crnd 1 and slot 1, previously seen accept (rnd 2) and responded with learn -> expected nil promise",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 1},
				},
				wantPrm: nil,
			},
		},
	},
	{
		NewAcceptor(),
		[]acceptorAction{
			{
				desc:    "prepare with crnd 2, no previous prepare or accepts -> expected promise with correct rnd and no slots",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 0},
					Crnd: &pb.Round{Id: 2},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 2},
				},
			},
		},
	},
	{
		NewAcceptor(),
		[]acceptorAction{ // TODO(meling) these two tests does not make sense to me; the expected output is a promise not prepare; check
			{
				desc:    "prepare for slot 1 with round 1, no previous history (slots) -> expected prepare with correct rnd and not slots",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Crnd: &pb.Round{Id: 1},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 1},
				},
			},
			{
				desc:    "another prepare for slot 1 with round 2 -> expected another prepare with correct rnd and not slots",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Crnd: &pb.Round{Id: 2},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 2},
				},
			},
		},
	},
	{
		NewAcceptor(),
		[]acceptorAction{
			{
				desc:    "prepare for slot 1, crnd 2 with no previous history -> expected correct promise",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 2},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 2},
				},
			},
			{
				desc:    "accept for slot 1 with current round -> expected learn with correct slot, rnd and value",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
			},
			{
				desc:    "new prepare for slot 1 with higher crnd -> expected promise with correct rnd and history (slot 1)",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 3},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 3},
					Slots: []*pb.PromiseSlot{
						{
							Slot:  &pb.Slot{Id: 1},
							Vrnd:  &pb.Round{Id: 2},
							Value: &testingValueOne,
						},
					},
				},
			},
		},
	},
	{
		NewAcceptor(),
		[]acceptorAction{
			{
				desc:    "prepare for slot 1, crnd 2, no previous history -> expected correct promise",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 2},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 2},
				},
			},
			{
				desc:    "accept for slot 1 with current rnd -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
			},
			{
				desc:    "accept for slot 1 with _higher_ rnd -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 5},
					Val:  &testingValueTwo,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 5},
					Val:  &testingValueTwo,
				},
			},
			{
				desc:    "accept for slot 2 with _higher_ rnd -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 2},
					Rnd:  &pb.Round{Id: 8},
					Val:  &testingValueThree,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 2},
					Rnd:  &pb.Round{Id: 8},
					Val:  &testingValueThree,
				},
			},
		},
	},
	{
		NewAcceptor(),
		[]acceptorAction{
			{
				desc:    "prepare for slot 1, crnd 4, no previous history -> expected correct promise",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 4},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 4},
				},
			},
			{
				desc:    "accept with _lower_ rnd (2) than out current rnd (4) -> expected nil learn, ignoring the accept",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
				wantLrn: nil,
			},
		},
	},
	{
		NewAcceptor(),
		[]acceptorAction{
			{
				desc:    "prepare for slot 1, crnd 2, no previous history -> expected correct promise",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 2},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 2},
				},
			},
			{
				desc:    "accept for slot 1, rnd 2 -> expected learn with correct slot, rnd and value",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
			},
			{
				desc:    "accept for slot 2, rnd 2 -> expected learn with correct slot, rnd and value",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 3},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueTwo,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 3},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueTwo,
				},
			},
			{
				desc:    "accept for slot 4, rnd 2 -> expected learn with correct slot, rnd and value",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 4},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueThree,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 4},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueThree,
				},
			},
			{
				desc:    "new prepare for slot 2, crnd 3 -> expected promise with correct rnd and history (slot 3 and 4)",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 2},
					Crnd: &pb.Round{Id: 3},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 3},
					Slots: []*pb.PromiseSlot{
						{
							Slot:  &pb.Slot{Id: 3},
							Vrnd:  &pb.Round{Id: 2},
							Value: &testingValueTwo,
						},
						{
							Slot:  &pb.Slot{Id: 4},
							Vrnd:  &pb.Round{Id: 2},
							Value: &testingValueThree,
						},
					},
				},
			},
		},
	},
	{
		NewAcceptor(),
		[]acceptorAction{
			{
				desc:    "prepare for slot 1, crnd 2, no previous history -> expected correct promise",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 2},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 2},
				},
			},
			{
				desc:    "accept for slot 1, rnd 2 -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
			},
			{
				desc:    "accept for slot 3, with different sender and lower round -> expected nil learn, ignoring the accept",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 3},
					Rnd:  &pb.Round{Id: 1},
					Val:  &testingValueTwo,
				},
				wantLrn: nil,
			},
			{
				desc:    "accept for slot 4, rnd 2 -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 4},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueThree,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 4},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueThree,
				},
			},
			{
				desc:    "accept for slot 5 with _higher_ rnd (5) -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 5},
					Rnd:  &pb.Round{Id: 5},
					Val:  &testingValueOne,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 5},
					Rnd:  &pb.Round{Id: 5},
					Val:  &testingValueOne,
				},
			},
			{
				desc:    "accept for slot 7 with old lower rnd (1) -> expected nil learn, ignoring the accept",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 7},
					Rnd:  &pb.Round{Id: 1},
					Val:  &testingValueTwo,
				},
				wantLrn: nil,
			},
			{
				desc:    "new prepare for slot 2 and higher round (7) -> expected promise with correct rnd (7) and history (slot 4 and 5)",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 2},
					Crnd: &pb.Round{Id: 7},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 7},
					Slots: []*pb.PromiseSlot{
						{
							Slot:  &pb.Slot{Id: 4},
							Vrnd:  &pb.Round{Id: 2},
							Value: &testingValueThree,
						},
						{
							Slot:  &pb.Slot{Id: 5},
							Vrnd:  &pb.Round{Id: 5},
							Value: &testingValueOne,
						},
					},
				},
			},
		},
	},
	{
		NewAcceptor(),
		[]acceptorAction{
			{
				desc:    "prepare for slot 1, crnd 2, no previous history -> expected correct promise",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 2},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 2},
				},
			},
			{
				desc:    "accept for slot 1, rnd 2 -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
			},
			{
				desc:    "accept for slot 3, with different sender and lower round -> expected nil learn, ignoring the accept",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 3},
					Rnd:  &pb.Round{Id: 1},
					Val:  &testingValueTwo,
				},
				wantLrn: nil,
			},
			{
				desc:    "accept for slot 4, rnd 2 -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 4},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueThree,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 4},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueThree,
				},
			},
			{
				desc:    "accept for slot 6 with _higher_ rnd (5) -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 6},
					Rnd:  &pb.Round{Id: 5},
					Val:  &testingValueOne,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 6},
					Rnd:  &pb.Round{Id: 5},
					Val:  &testingValueOne,
				},
			},
			{
				desc:    "accept for slot 7 with old lower rnd (1) -> expected nil learn, ignoring the accept",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 7},
					Rnd:  &pb.Round{Id: 1},
					Val:  &testingValueTwo,
				},
				wantLrn: nil,
			},
			{
				desc:    "new prepare for slot 2 and higher round (7) -> expected promise with correct rnd (7) and history (slot 4 and 6)",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 2},
					Crnd: &pb.Round{Id: 7},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 7},
					Slots: []*pb.PromiseSlot{
						{
							Slot:  &pb.Slot{Id: 4},
							Vrnd:  &pb.Round{Id: 2},
							Value: &testingValueThree,
						},
						{
							Slot:  &pb.Slot{Id: 6},
							Vrnd:  &pb.Round{Id: 5},
							Value: &testingValueOne,
						},
					},
				},
			},
		},
	},
	{
		NewAcceptor(),
		[]acceptorAction{
			{
				desc:    "prepare for slot 1, crnd 2, no previous history -> expected correct promise",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 2},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 2},
				},
			},
			{
				desc:    "accept for slot 1, rnd 2 -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
			},
			{
				desc:    "accept for slot 1 again but with rnd 5 -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 5},
					Val:  &testingValueOne,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 5},
					Val:  &testingValueOne,
				},
			},
			{
				desc:    "accept for slot 3, with different sender and lower round -> expected nil learn, ignoring the accept",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 3},
					Rnd:  &pb.Round{Id: 1},
					Val:  &testingValueTwo,
				},
				wantLrn: nil,
			},
			{
				desc:    "accept for slot 4, rnd 5 -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 4},
					Rnd:  &pb.Round{Id: 5},
					Val:  &testingValueThree,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 4},
					Rnd:  &pb.Round{Id: 5},
					Val:  &testingValueThree,
				},
			},
			{
				desc:    "accept again for slot 4, with _higher_ round, rnd 8 -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 4},
					Rnd:  &pb.Round{Id: 8},
					Val:  &testingValueThree,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 4},
					Rnd:  &pb.Round{Id: 8},
					Val:  &testingValueThree,
				},
			},
			{
				desc:    "accept for slot 6 with _higher_ rnd (11) -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 6},
					Rnd:  &pb.Round{Id: 11},
					Val:  &testingValueOne,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 6},
					Rnd:  &pb.Round{Id: 11},
					Val:  &testingValueOne,
				},
			},
			{
				desc:    "accept for slot 7 with old lower rnd (1) -> expected nil learn, ignoring the accept",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 7},
					Rnd:  &pb.Round{Id: 1},
					Val:  &testingValueTwo,
				},
				wantLrn: nil,
			},
			{
				desc:    "new prepare for slot 2 and higher round (13) from 1 -> expected promise with correct rnd (13) and history (slot 4 and 6)",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 2},
					Crnd: &pb.Round{Id: 13},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 13},
					Slots: []*pb.PromiseSlot{
						{
							Slot:  &pb.Slot{Id: 4},
							Vrnd:  &pb.Round{Id: 8},
							Value: &testingValueThree,
						},
						{
							Slot:  &pb.Slot{Id: 6},
							Vrnd:  &pb.Round{Id: 11},
							Value: &testingValueOne,
						},
					},
				},
			},
		},
	},
	{
		NewAcceptor(),
		[]acceptorAction{
			{
				desc:    "prepare for slot 1, crnd 2, no previous history -> expected correct promise",
				msgtype: prepare,
				prepare: &pb.PrepareMsg{
					Slot: &pb.Slot{Id: 1},
					Crnd: &pb.Round{Id: 2},
				},
				wantPrm: &pb.PromiseMsg{
					Rnd: &pb.Round{Id: 2},
				},
			},
			{
				desc:    "accept for slot 1 with current rnd -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 2},
					Val:  &testingValueOne,
				},
			},
			{
				desc:    "accept for slot 2 with _higher_ rnd -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 2},
					Rnd:  &pb.Round{Id: 5},
					Val:  &testingValueThree,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 2},
					Rnd:  &pb.Round{Id: 5},
					Val:  &testingValueThree,
				},
			},
			{
				desc:    "accept for slot 1 again with _higher_ rnd after we previously have sent accept for slot 2 -> expected correct learn",
				msgtype: accept,
				accept: &pb.AcceptMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 8},
					Val:  &testingValueTwo,
				},
				wantLrn: &pb.LearnMsg{
					Slot: &pb.Slot{Id: 1},
					Rnd:  &pb.Round{Id: 8},
					Val:  &testingValueTwo,
				},
			},
		},
	},
}
