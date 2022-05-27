package gorumspaxos

import (
	"reflect"
	"testing"

	pb "dat520/lab5/gorumspaxos/proto"
)

func TestHandlePrepareQF(t *testing.T) {
	for _, test := range proposerQFTests {
		qspecs := NewPaxosQSpec(test.quorumSize)
		replies := make(map[uint32]*pb.PromiseMsg)
		for i, promiseMsg := range test.promiseMsgs {
			replies[uint32(i)] = promiseMsg
		}
		gotPSlots, gotOutput := qspecs.PrepareQF(test.inputPrepareMsg, replies)
		switch {
		case !test.wantOutput && gotOutput:
			t.Errorf("\nPromiseQF test failed\ndescription: %s\nwant no output\ngot %v",
				test.description, gotPSlots)
		case test.wantOutput && !gotOutput:
			t.Errorf("\nPromiseQF test failed\ndescription: %s\nwant output %v\ngot no output",
				test.description, test.outputPromiseMsg)
		case test.wantOutput && gotOutput:
			if !comparePromiseMsgs(gotPSlots, test.outputPromiseMsg) {
				t.Errorf("\nPromiseQF test failed\ndescription: %s\nwant output %v\ngot output %v\n",
					test.description, test.outputPromiseMsg, gotPSlots)
			}
		}
	}
}

func TestHandleAcceptQF(t *testing.T) {
	for _, test := range acceptQFTests {
		qspecs := NewPaxosQSpec(test.quorumSize)
		replies := make(map[uint32]*pb.LearnMsg)
		for i, learnMsg := range test.learnMsgs {
			replies[uint32(i)] = learnMsg
		}
		gotLearnMsg, gotOutput := qspecs.AcceptQF(test.inputAcceptMsg, replies)
		switch {
		case !test.wantOutput && gotOutput:
			t.Errorf("\nPromiseQF test failed\ndescription: %s\nwant no output\ngot %v",
				test.description, gotLearnMsg)
		case test.wantOutput && !gotOutput:
			t.Errorf("\nPromiseQF test failed\ndescription: %s\nwant output %v\ngot no output",
				test.description, test.outputLearnMsg)
		case test.wantOutput && gotOutput:
			if !reflect.DeepEqual(gotLearnMsg, test.outputLearnMsg) {
				t.Errorf("\nPromiseQF test failed\ndescription: %s\nwant output %v\ngot output %v\n",
					test.description, test.outputLearnMsg, gotLearnMsg)
			}
		}
	}
}

func TestHandleClientReqQF(t *testing.T) {
	for _, test := range clientHandleQFTests {
		qspecs := NewPaxosQSpec(test.QuorumSize)
		replies := make(map[uint32]*pb.Response)
		for idx, reply := range test.replies {
			replies[uint32(idx)] = reply
		}
		resp, gotOutput := qspecs.ClientHandleQF(test.inputValue, replies)
		switch {
		case !test.wantOutput && gotOutput:
			t.Errorf("\nClientHandleQF test failed\ndescription: %s\nwant no output\ngot %v",
				test.description, resp)
		case test.wantOutput && !gotOutput:
			t.Errorf("\nClientHandleQF test failed\ndescription: %s\nwant output, got no output\n",
				test.description)
		}
	}
}

func TestHandleCommitQF(t *testing.T) {
	lrnMsg := &pb.LearnMsg{}
	for _, test := range commitQFTests {
		qspecs := NewPaxosQSpec(test.quorumSize)
		replies := make(map[uint32]*pb.Empty)
		for i := 0; i < test.testQuorumSize; i++ {
			replies[uint32(i)] = &pb.Empty{}
		}
		_, gotOutput := qspecs.CommitQF(lrnMsg, replies)
		switch {
		case !test.wantOutput && gotOutput:
			t.Errorf("\nCommitQF test failed\ndescription: %s\nwant no output\ngot some output",
				test.description)
		case test.wantOutput && !gotOutput:
			t.Errorf("\nCommitQF test failed\ndescription: %s\nwant output, got no output\n",
				test.description)
		}
	}
}

var proposerQFTests = []proposeQFTest{
	{
		description: "Valid prepare message, response with a valid promise message but no quorum",
		quorumSize:  2,
		promiseMsgs: []*pb.PromiseMsg{
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 1},
					},
				},
			},
		},
		inputPrepareMsg: &pb.PrepareMsg{
			Crnd: &pb.Round{Id: 6},
			Slot: &pb.Slot{Id: 2},
		},
		wantOutput:       false,
		outputPromiseMsg: &pb.PromiseMsg{},
	},
	{
		description: `Valid prepare message, response with a valid quorum of promise messages, 
		but no new promise slots`,
		quorumSize: 2,
		promiseMsgs: []*pb.PromiseMsg{
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 2},
						Value: &pb.Value{
							ClientID:      "1",
							ClientSeq:     2,
							ClientCommand: "ls",
						},
					},
				},
			},
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 2},
					},
				},
			},
		},
		inputPrepareMsg: &pb.PrepareMsg{
			Crnd: &pb.Round{Id: 6},
			Slot: &pb.Slot{Id: 2},
		},
		wantOutput: true,
		outputPromiseMsg: &pb.PromiseMsg{
			Rnd:   &pb.Round{Id: 6},
			Slots: make([]*pb.PromiseSlot, 0),
		},
	},
	{
		description: `Valid prepare message, response with a valid quorum of  some invalid promise messages, 
		but no new promise slots`,
		quorumSize: 2,
		promiseMsgs: []*pb.PromiseMsg{
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 1},
					},
				},
			},
			{
				Rnd: &pb.Round{Id: 3},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 2},
					},
				},
			},
		},
		inputPrepareMsg: &pb.PrepareMsg{
			Crnd: &pb.Round{Id: 6},
			Slot: &pb.Slot{Id: 2},
		},
		wantOutput: false,
	},
	{
		description: `Valid prepare message, response with a valid quorum of promise messages, 
		with new promise slots`,
		quorumSize: 2,
		promiseMsgs: []*pb.PromiseMsg{
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 2},
					},
				},
			},
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 1},
					},
				},
			},
		},
		inputPrepareMsg: &pb.PrepareMsg{
			Crnd: &pb.Round{Id: 6},
			Slot: &pb.Slot{Id: 1},
		},
		wantOutput: true,
		outputPromiseMsg: &pb.PromiseMsg{
			Rnd: &pb.Round{Id: 6},
			Slots: []*pb.PromiseSlot{
				{
					Slot: &pb.Slot{Id: 2},
				},
			},
		},
	},
	{
		description: `Valid prepare message, response with a valid quorum of promise messages, 
		with new promise slots check for sorted order`,
		quorumSize: 2,
		promiseMsgs: []*pb.PromiseMsg{
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 2},
					},
					{
						Slot: &pb.Slot{Id: 4},
					},
				},
			},
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 1},
					},
					{
						Slot: &pb.Slot{Id: 4},
					},
				},
			},
		},
		inputPrepareMsg: &pb.PrepareMsg{
			Crnd: &pb.Round{Id: 6},
			Slot: &pb.Slot{Id: 1},
		},
		wantOutput: true,
		outputPromiseMsg: &pb.PromiseMsg{
			Rnd: &pb.Round{Id: 6},
			Slots: []*pb.PromiseSlot{
				{
					Slot: &pb.Slot{Id: 2},
				},
				{
					Slot:  &pb.Slot{Id: 3},
					Vrnd:  &pb.Round{Id: 6},
					Value: &pb.Value{IsNoop: true},
				},
				{
					Slot: &pb.Slot{Id: 4},
				},
			},
		},
	},
	{
		description: `Valid prepare message, response with a valid quorum of promise messages, 
		with new promise slots with no op flags because of gaps `,
		quorumSize: 2,
		promiseMsgs: []*pb.PromiseMsg{
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 2},
					},
					{
						Slot: &pb.Slot{Id: 4},
					},
					{
						Slot: &pb.Slot{Id: 6},
					},
				},
			},
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 1},
					},
					{
						Slot: &pb.Slot{Id: 3},
					},
				},
			},
		},
		inputPrepareMsg: &pb.PrepareMsg{
			Crnd: &pb.Round{Id: 6},
			Slot: &pb.Slot{Id: 1},
		},
		wantOutput: true,
		outputPromiseMsg: &pb.PromiseMsg{
			Rnd: &pb.Round{Id: 6},
			Slots: []*pb.PromiseSlot{
				{
					Slot: &pb.Slot{Id: 2},
				},
				{
					Slot: &pb.Slot{Id: 3},
				},
				{
					Slot: &pb.Slot{Id: 4},
				},
				{
					Slot:  &pb.Slot{Id: 5},
					Value: &pb.Value{IsNoop: true},
					Vrnd:  &pb.Round{Id: 6},
				},
				{
					Slot: &pb.Slot{Id: 6},
				},
			},
		},
	},
	{
		description: `Valid prepare message, response with a valid quorum of promise messages, 
		with new promise slots with no op flags because of gaps and round of accept should be proper `,
		quorumSize: 2,
		promiseMsgs: []*pb.PromiseMsg{
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 2},
						Vrnd: &pb.Round{Id: 3},
					},
					{
						Slot: &pb.Slot{Id: 3},
						Vrnd: &pb.Round{Id: 4},
					},
					{
						Slot: &pb.Slot{Id: 6},
					},
				},
			},
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 2},
						Vrnd: &pb.Round{Id: 4},
					},
					{
						Slot: &pb.Slot{Id: 3},
						Vrnd: &pb.Round{Id: 3},
					},
				},
			},
		},
		inputPrepareMsg: &pb.PrepareMsg{
			Crnd: &pb.Round{Id: 6},
			Slot: &pb.Slot{Id: 1},
		},
		wantOutput: true,
		outputPromiseMsg: &pb.PromiseMsg{
			Rnd: &pb.Round{Id: 6},
			Slots: []*pb.PromiseSlot{
				{
					Slot: &pb.Slot{Id: 2},
					Vrnd: &pb.Round{Id: 4},
				},
				{
					Slot: &pb.Slot{Id: 3},
					Vrnd: &pb.Round{Id: 4},
				},
				{
					Slot:  &pb.Slot{Id: 4},
					Value: &pb.Value{IsNoop: true},
					Vrnd:  &pb.Round{Id: 6},
				},
				{
					Slot:  &pb.Slot{Id: 5},
					Value: &pb.Value{IsNoop: true},
					Vrnd:  &pb.Round{Id: 6},
				},
				{
					Slot: &pb.Slot{Id: 6},
				},
			},
		},
	},
	{
		description: `Valid prepare message, response with a valid quorum of promise messages, 
		with  one invalid promise message `,
		quorumSize: 2,
		promiseMsgs: []*pb.PromiseMsg{
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 2},
						Vrnd: &pb.Round{Id: 3},
					},
					{
						Slot: &pb.Slot{Id: 3},
						Vrnd: &pb.Round{Id: 4},
					},
					{
						Slot: &pb.Slot{Id: 6},
					},
				},
			},
			{
				// This message should not be processed as round variable
				// is not available
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 2},
						Vrnd: &pb.Round{Id: 6},
					},
					{
						Slot: &pb.Slot{Id: 8},
						Vrnd: &pb.Round{Id: 4},
					},
					{
						Slot: &pb.Slot{Id: 7},
					},
				},
			},
			{
				Rnd: &pb.Round{Id: 6},
				Slots: []*pb.PromiseSlot{
					{
						Slot: &pb.Slot{Id: 2},
						Vrnd: &pb.Round{Id: 4},
					},
					{
						Slot: &pb.Slot{Id: 3},
						Vrnd: &pb.Round{Id: 3},
					},
				},
			},
		},
		inputPrepareMsg: &pb.PrepareMsg{
			Crnd: &pb.Round{Id: 6},
			Slot: &pb.Slot{Id: 1},
		},
		wantOutput: true,
		outputPromiseMsg: &pb.PromiseMsg{
			Rnd: &pb.Round{Id: 6},
			Slots: []*pb.PromiseSlot{
				{
					Slot: &pb.Slot{Id: 2},
					Vrnd: &pb.Round{Id: 4},
				},
				{
					Slot: &pb.Slot{Id: 3},
					Vrnd: &pb.Round{Id: 4},
				},
				{
					Slot:  &pb.Slot{Id: 4},
					Value: &pb.Value{IsNoop: true},
					Vrnd:  &pb.Round{Id: 6},
				},
				{
					Slot:  &pb.Slot{Id: 5},
					Value: &pb.Value{IsNoop: true},
					Vrnd:  &pb.Round{Id: 6},
				},
				{
					Slot: &pb.Slot{Id: 6},
				},
			},
		},
	},
}

var acceptQFTests = []acceptQFTest{
	{
		description: "No quorum of learn messages implies no output",
		quorumSize:  2,
		learnMsgs: []*pb.LearnMsg{
			{
				Rnd:  &pb.Round{},
				Val:  &pb.Value{},
				Slot: &pb.Slot{},
			},
		},
		inputAcceptMsg: &pb.AcceptMsg{
			Rnd:  &pb.Round{},
			Val:  &pb.Value{},
			Slot: &pb.Slot{},
		},
		wantOutput: false,
	},
	{
		description: "Quorum of learn messages with some invalid learn messages causing no output ",
		quorumSize:  2,
		learnMsgs: []*pb.LearnMsg{
			{
				Rnd: &pb.Round{Id: 3},
				Val: &pb.Value{
					ClientID:      "1",
					ClientSeq:     1,
					ClientCommand: "rm",
				},
				Slot: &pb.Slot{Id: 1},
			},
			{
				Rnd: &pb.Round{Id: 3},
				Val: &pb.Value{
					ClientID:      "1",
					ClientSeq:     1,
					ClientCommand: "rm",
				},
				Slot: &pb.Slot{Id: 2},
			},
		},
		inputAcceptMsg: &pb.AcceptMsg{
			Rnd:  &pb.Round{Id: 3},
			Slot: &pb.Slot{Id: 2},
			Val: &pb.Value{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
		},
		wantOutput: false,
	},
	{
		description: `Quorum of valid learn messages, but round learn messages 
		is less than input accept message implies no output`,
		quorumSize: 2,
		learnMsgs: []*pb.LearnMsg{
			{
				Rnd: &pb.Round{Id: 2},
				Val: &pb.Value{
					ClientID:      "1",
					ClientSeq:     1,
					ClientCommand: "rm",
				},
				Slot: &pb.Slot{Id: 2},
			},
			{
				Rnd: &pb.Round{Id: 3},
				Val: &pb.Value{
					ClientID:      "1",
					ClientSeq:     1,
					ClientCommand: "rm",
				},
				Slot: &pb.Slot{Id: 2},
			},
		},
		inputAcceptMsg: &pb.AcceptMsg{
			Rnd:  &pb.Round{Id: 3},
			Slot: &pb.Slot{Id: 2},
			Val: &pb.Value{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
		},
		wantOutput: false,
	},
	{
		description: `Quorum of valid learn messages with some learn messages having 
		different value than input accept messase, implies no output`,
		quorumSize: 2,
		learnMsgs: []*pb.LearnMsg{
			{
				Rnd: &pb.Round{Id: 3},
				Val: &pb.Value{
					ClientID:      "1",
					ClientSeq:     1,
					ClientCommand: "rm",
				},
				Slot: &pb.Slot{Id: 2},
			},
			{
				Rnd: &pb.Round{Id: 3},
				Val: &pb.Value{
					ClientID:      "1",
					ClientSeq:     2,
					ClientCommand: "rm",
				},
				Slot: &pb.Slot{Id: 2},
			},
		},
		inputAcceptMsg: &pb.AcceptMsg{
			Rnd:  &pb.Round{Id: 3},
			Slot: &pb.Slot{Id: 2},
			Val: &pb.Value{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
		},
		wantOutput: false,
	},
	{
		description: `Quorum of valid learn messages with a quorum of learn messages 
		having same round as input accept message implies output `,
		quorumSize: 2,
		learnMsgs: []*pb.LearnMsg{
			{
				Rnd: &pb.Round{Id: 3},
				Val: &pb.Value{
					ClientID:      "1",
					ClientSeq:     1,
					ClientCommand: "rm",
				},
				Slot: &pb.Slot{Id: 2},
			},
			{
				Rnd: &pb.Round{Id: 3},
				Val: &pb.Value{
					ClientID:      "1",
					ClientSeq:     1,
					ClientCommand: "rm",
				},
				Slot: &pb.Slot{Id: 2},
			},
		},
		inputAcceptMsg: &pb.AcceptMsg{
			Rnd:  &pb.Round{Id: 3},
			Slot: &pb.Slot{Id: 2},
			Val: &pb.Value{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
		},
		outputLearnMsg: &pb.LearnMsg{
			Rnd:  &pb.Round{Id: 3},
			Slot: &pb.Slot{Id: 2},
			Val: &pb.Value{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
		},
		wantOutput: true,
	},
}

var clientHandleQFTests = []clientHandleQFTest{
	{
		description: "quorum size valid replies",
		inputValue:  &pb.Value{ClientID: "1", ClientSeq: 1, ClientCommand: "rm"},
		QuorumSize:  3,
		replies: []*pb.Response{
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
		},
		wantOutput: true,
	},
	{
		description: "quorum size replies with invalid ClientSeq in one reply",
		inputValue:  &pb.Value{ClientID: "1", ClientSeq: 1, ClientCommand: "rm"},
		QuorumSize:  3,
		replies: []*pb.Response{
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     2,
				ClientCommand: "rm",
			},
		},
		wantOutput: false,
	},
	{
		description: "quorum size replies with invalid ClientCommand in one reply",
		inputValue:  &pb.Value{ClientID: "1", ClientSeq: 1, ClientCommand: "rm"},
		QuorumSize:  3,
		replies: []*pb.Response{
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "1rm",
			},
		},
		wantOutput: false,
	},
	{
		description: "quorum size replies with invalid ClientId in one reply",
		inputValue:  &pb.Value{ClientID: "1", ClientSeq: 1, ClientCommand: "rm"},
		QuorumSize:  3,
		replies: []*pb.Response{
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "2",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
		},
		wantOutput: false,
	},
	{
		description: "All replies with invalid ClientSeq in one reply, still valid quorum",
		inputValue:  &pb.Value{ClientID: "1", ClientSeq: 1, ClientCommand: "rm"},
		QuorumSize:  2,
		replies: []*pb.Response{
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     2,
				ClientCommand: "rm",
			},
		},
		wantOutput: true,
	},
	{
		description: "All replies with invalid ClientCommand in one reply, still valid quorum",
		inputValue:  &pb.Value{ClientID: "1", ClientSeq: 1, ClientCommand: "rm"},
		QuorumSize:  2,
		replies: []*pb.Response{
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "1rm",
			},
		},
		wantOutput: true,
	},
	{
		description: "All replies with invalid ClientId in one reply, still valid quorum",
		inputValue:  &pb.Value{ClientID: "1", ClientSeq: 1, ClientCommand: "rm"},
		QuorumSize:  2,
		replies: []*pb.Response{
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "2",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
		},
		wantOutput: true,
	},
	{
		description: "No quorum replies",
		inputValue:  &pb.Value{ClientID: "1", ClientSeq: 1, ClientCommand: "rm"},
		QuorumSize:  3,
		replies: []*pb.Response{
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
			{
				ClientID:      "1",
				ClientSeq:     1,
				ClientCommand: "rm",
			},
		},
		wantOutput: false,
	},
}

var commitQFTests = []commitQFTest{
	{
		description:    "quorum size replies",
		quorumSize:     3,
		testQuorumSize: 3,
		wantOutput:     true,
	},
	{
		description:    "no quorum size replies, so no output",
		quorumSize:     3,
		testQuorumSize: 2,
		wantOutput:     false,
	},
}
